package cli

import (
	"errors"
	"fmt"

	"github.com/Blazzical/snugNAS/internal/compose"
	"github.com/Blazzical/snugNAS/internal/config"
	"github.com/Blazzical/snugNAS/internal/installer"
	"github.com/spf13/cobra"
)

func doctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Full health check — prerequisites, config, container status",
		Long: `Combines snugnas install + snugnas status into one report. Useful
when something's not working and you want one place to see whether the
fault is in Docker, the config, or a specific container.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			ok := true

			fmt.Fprintln(out, "== Prerequisites ==")
			s := installer.Check(cmd.Context())
			mark := func(b bool) string {
				if b {
					return "[ok]"
				}
				ok = false
				return "[--]"
			}
			fmt.Fprintln(out, mark(s.DockerInstalled), "Docker installed")
			fmt.Fprintln(out, mark(s.DockerRunning), "Docker daemon reachable")
			fmt.Fprintln(out, mark(s.WSL2Default), "WSL2 default version")
			for _, n := range s.Notes {
				fmt.Fprintln(out, "    *", n)
			}

			fmt.Fprintln(out, "\n== Config ==")
			cfg, err := config.Load()
			if err != nil {
				if config.IsNotConfigured(err) {
					p, _ := config.Path()
					fmt.Fprintln(out, "[--] no config at", p)
					fmt.Fprintln(out, "     run `snugnas wizard` to create one")
					ok = false
				} else {
					fmt.Fprintln(out, "[--] error loading config:", err)
					ok = false
				}
			} else {
				p, _ := config.Path()
				fmt.Fprintln(out, "[ok] config loaded from", p)
				fmt.Fprintln(out, "     storage_root:", cfg.StorageRoot)
				fmt.Fprintln(out, "     hostname:    ", cfg.Hostname)
				fmt.Fprintln(out, "     dashboard:   ", fmt.Sprintf("http://127.0.0.1:%d/", cfg.DashboardPort))
			}

			fmt.Fprintln(out, "\n== Services ==")
			if cfg == nil {
				fmt.Fprintln(out, "[--] skipped — config not loaded")
				ok = false
			} else if !s.DockerRunning {
				fmt.Fprintln(out, "[--] skipped — Docker not running")
			} else {
				statuses, err := compose.PS(cmd.Context())
				if err != nil {
					fmt.Fprintln(out, "[--] docker compose ps failed:", err)
					ok = false
				} else if len(statuses) == 0 {
					fmt.Fprintln(out, "[--] no containers running — try `snugnas up`")
					ok = false
				} else {
					for _, svc := range statuses {
						state := svc.State
						if svc.Health != "" {
							state += " / " + svc.Health
						}
						label := "[ok]"
						if svc.State != "running" || svc.Health == "unhealthy" {
							label = "[--]"
							ok = false
						}
						fmt.Fprintf(out, "%s %-26s %s\n", label, svc.Service, state)
					}
				}
			}

			fmt.Fprintln(out)
			if !ok {
				return errors.New("doctor found issues")
			}
			fmt.Fprintln(out, "All checks passed.")
			return nil
		},
	}
}
