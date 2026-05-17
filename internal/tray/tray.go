//go:build windows

// Package tray runs the snugNAS Windows tray app. It owns the lifecycle of the
// dashboard HTTP server and the mDNS publisher so the user has a single
// always-on entry point: launch snugNAS, see the tray icon, click to open the
// dashboard. Backed by github.com/getlantern/systray, which uses Win32 syscalls
// on Windows (no Cgo). The non-Windows build sits in tray_other.go.
package tray

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"sync"

	"github.com/Blazzical/snugNAS/internal/config"
	"github.com/Blazzical/snugNAS/internal/landing"
	"github.com/Blazzical/snugNAS/internal/mdns"
	"github.com/getlantern/systray"
)

// Run starts the tray app. Must be called on the main goroutine — systray.Run
// requires it.
func Run(cfg *config.Config) error {
	t := &tray{cfg: cfg}
	systray.Run(t.onReady, t.onExit)
	if t.serverErr != nil && !errors.Is(t.serverErr, context.Canceled) {
		return t.serverErr
	}
	return nil
}

type tray struct {
	cfg       *config.Config
	cancel    context.CancelFunc
	pub       *mdns.Publisher
	wg        sync.WaitGroup
	serverErr error
}

func (t *tray) onReady() {
	systray.SetIcon(snugnasIcon())
	systray.SetTitle("snugNAS")
	systray.SetTooltip(fmt.Sprintf("snugNAS — http://%s:8080/", t.cfg.Hostname))

	mOpen := systray.AddMenuItem("Open dashboard", "")
	mOpenLAN := systray.AddMenuItem(
		fmt.Sprintf("Open on LAN (http://%s:8080/)", t.cfg.Hostname),
		"Opens via the mDNS hostname — works from this PC if Bonjour is installed",
	)
	systray.AddSeparator()
	mStatus := systray.AddMenuItem("Status: starting…", "")
	mStatus.Disable()
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit snugNAS", "Stop the dashboard and exit")

	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	if pub, err := mdns.Publish(t.cfg.Hostname, t.cfg.DashboardPort); err != nil {
		mStatus.SetTitle("Status: mDNS failed — dashboard only")
	} else {
		t.pub = pub
		mStatus.SetTitle(fmt.Sprintf("Status: advertising %s", pub.Hostname))
	}

	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		if err := landing.Serve(ctx, t.cfg); err != nil && !errors.Is(err, context.Canceled) {
			t.serverErr = err
		}
	}()

	go func() {
		dashURL := fmt.Sprintf("http://127.0.0.1:%d/", t.cfg.DashboardPort)
		lanURL := fmt.Sprintf("http://%s:8080/", t.cfg.Hostname)
		for {
			select {
			case <-mOpen.ClickedCh:
				_ = openURL(dashURL)
			case <-mOpenLAN.ClickedCh:
				_ = openURL(lanURL)
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func (t *tray) onExit() {
	if t.cancel != nil {
		t.cancel()
	}
	if t.pub != nil {
		t.pub.Shutdown()
	}
	t.wg.Wait()
}

func openURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
