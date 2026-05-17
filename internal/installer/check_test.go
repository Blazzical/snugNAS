package installer

import "testing"

func TestStatusOK(t *testing.T) {
	cases := []struct {
		name string
		s    Status
		want bool
	}{
		{"empty", Status{}, false},
		{"docker installed only", Status{DockerInstalled: true}, false},
		{"docker installed + running", Status{DockerInstalled: true, DockerRunning: true}, false},
		{"all true", Status{DockerInstalled: true, DockerRunning: true, WSL2Default: true}, true},
		{"wsl only", Status{WSL2Default: true}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.s.OK(); got != c.want {
				t.Errorf("OK() = %v, want %v", got, c.want)
			}
		})
	}
}
