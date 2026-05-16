package cli

import (
	"errors"
	"fmt"

	"github.com/Blazzical/snugNAS/internal/installer"
	"github.com/spf13/cobra"
)

func installCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Check Docker Desktop and WSL2 are ready; print remediation if not",
		Long: `Probes the local machine for snugNAS prerequisites:
  - docker on PATH
  - Docker daemon reachable
  - WSL2 set as the default WSL version

Prints a checklist and exits non-zero if anything is missing, with a one-line
remediation hint for each failure. Does not install anything itself.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := installer.Check(cmd.Context())

			mark := func(b bool) string {
				if b {
					return "[ok]"
				}
				return "[--]"
			}
			out := cmd.OutOrStdout()
			fmt.Fprintln(out, mark(s.DockerInstalled), "Docker installed")
			fmt.Fprintln(out, mark(s.DockerRunning), "Docker daemon reachable")
			fmt.Fprintln(out, mark(s.WSL2Default), "WSL2 is the default version")

			if len(s.Notes) > 0 {
				fmt.Fprintln(out)
				for _, n := range s.Notes {
					fmt.Fprintln(out, "*", n)
				}
			}

			if !s.OK() {
				return errors.New("prerequisites not met")
			}
			fmt.Fprintln(out, "\nAll prerequisites met — snugnas wizard is ready to run.")
			return nil
		},
	}
}
