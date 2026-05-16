package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func installCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Detect or guide installation of Docker Desktop and WSL2",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO(installer): probe for docker.exe on PATH, check `wsl --status`,
			// guide the user through https://www.docker.com/products/docker-desktop/
			// if missing. See internal/installer.
			fmt.Println("install: not yet implemented")
			return nil
		},
	}
}
