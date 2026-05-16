// Package wizard runs the first-run web wizard. The Go binary serves a small
// single-page form on 127.0.0.1:<port>; the user fills in storage path,
// hostname, and admin email; submission writes config.toml and kicks off
// `docker compose up` in a goroutine. The page polls /api/state for progress.
package wizard

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Blazzical/snugNAS/internal/config"
	"github.com/Blazzical/snugNAS/internal/landing"
)

//go:embed web
var webFS embed.FS

// Phase is the lifecycle of the wizard run.
type Phase string

const (
	PhaseForm         Phase = "form"
	PhaseProvisioning Phase = "provisioning"
	PhaseDone         Phase = "done"
	PhaseFailed       Phase = "failed"
)

// State is the wizard's externally-observable state, returned as JSON to the
// polling client.
type State struct {
	Phase Phase    `json:"phase"`
	Logs  []string `json:"logs"`
	Error string   `json:"error,omitempty"`
	// PostUpURL is where the browser should redirect on success.
	PostUpURL string `json:"post_up_url,omitempty"`
}

// FormData is the payload the wizard page POSTs to /api/commit.
type FormData struct {
	StorageRoot string `json:"storage_root"`
	Hostname    string `json:"hostname"`
	AdminEmail  string `json:"admin_email"`
}

// Server holds the in-memory wizard state.
type Server struct {
	mu    sync.Mutex
	phase Phase
	logs  []string
	err   error
	cfg   *config.Config
}

// Serve runs the wizard HTTP server on 127.0.0.1:port and blocks until ctx is
// canceled.
func Serve(ctx context.Context, port int) error {
	s := &Server{phase: PhaseForm}

	// Pre-populate form defaults from any existing config so a re-run shows
	// the user's previous answers.
	if existing, err := config.Load(); err == nil {
		s.cfg = existing
	} else {
		s.cfg = config.Default()
	}

	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		return err
	}

	// Wizard routes — active during PhaseForm and PhaseProvisioning.
	wizardMux := http.NewServeMux()
	wizardMux.Handle("/static/", http.FileServer(http.FS(sub)))
	wizardMux.HandleFunc("/", s.handleIndex(sub))

	// API endpoints are always served by the wizard, regardless of phase,
	// so the page can keep polling after provisioning completes.
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/defaults", s.handleDefaults)
	apiMux.HandleFunc("/api/commit", s.handleCommit)
	apiMux.HandleFunc("/api/state", s.handleState)

	// Dashboard routes — activated after PhaseDone so /, /static/*, and
	// /api/qr serve the snugNAS landing page instead of the wizard form.
	dashboardMux := landing.Handlers(s.cfg)

	// Phased router: API always goes to wizard; everything else goes to
	// wizard normally and to dashboard once provisioning is done.
	phased := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/defaults" || r.URL.Path == "/api/commit" || r.URL.Path == "/api/state" {
			apiMux.ServeHTTP(w, r)
			return
		}
		s.mu.Lock()
		done := s.phase == PhaseDone
		s.mu.Unlock()
		if done {
			dashboardMux.ServeHTTP(w, r)
			return
		}
		wizardMux.ServeHTTP(w, r)
	})

	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	srv := &http.Server{
		Addr:              addr,
		Handler:           phased,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case <-ctx.Done():
		sCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(sCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func (s *Server) handleIndex(sub fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		data, err := fs.ReadFile(sub, "index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(data)
	}
}

func (s *Server) handleDefaults(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	writeJSON(w, http.StatusOK, FormData{
		StorageRoot: s.cfg.StorageRoot,
		Hostname:    s.cfg.Hostname,
		AdminEmail:  s.cfg.AdminEmail,
	})
}

func (s *Server) handleCommit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var form FormData
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := validateForm(form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	if s.phase == PhaseProvisioning {
		s.mu.Unlock()
		http.Error(w, "provisioning already in progress", http.StatusConflict)
		return
	}
	s.cfg.StorageRoot = form.StorageRoot
	s.cfg.Hostname = form.Hostname
	s.cfg.AdminEmail = form.AdminEmail
	s.phase = PhaseProvisioning
	s.logs = nil
	s.err = nil
	s.mu.Unlock()

	go s.provision()

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	st := State{
		Phase: s.phase,
		Logs:  append([]string(nil), s.logs...),
	}
	if s.err != nil {
		st.Error = s.err.Error()
	}
	if s.phase == PhaseDone {
		st.PostUpURL = "http://localhost:8080/"
	}
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, st)
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func validateForm(f FormData) error {
	if f.StorageRoot == "" {
		return errors.New("storage_root is required")
	}
	if f.Hostname == "" {
		return errors.New("hostname is required")
	}
	return nil
}
