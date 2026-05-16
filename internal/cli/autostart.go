package cli

import (
	"fmt"

	"github.com/Blazzical/snugNAS/internal/autostart"
	"github.com/spf13/cobra"
)

func autostartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autostart",
		Short: "Run snugNAS automatically when you log in (Windows only)",
		Long: `Manages an entry in HKCU\Software\Microsoft\Windows\CurrentVersion\Run
that launches "snugnas tray" at login. Effective for the current user only —
no admin rights required.`,
	}
	cmd.AddCommand(
		autostartEnableCmd(),
		autostartDisableCmd(),
		autostartStatusCmd(),
	)
	return cmd
}

func autostartEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable",
		Short: "Enable autostart at login",
		RunE: func(cmd *cobra.Command, args []string) error {
			line, err := autostart.Enable()
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Autostart enabled. On next login the following will run:")
			fmt.Fprintln(cmd.OutOrStdout(), " ", line)
			return nil
		},
	}
}

func autostartDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Short: "Disable autostart at login",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := autostart.Disable(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Autostart disabled.")
			return nil
		},
	}
}

func autostartStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Print whether autostart is currently enabled",
		RunE: func(cmd *cobra.Command, args []string) error {
			enabled, line, err := autostart.Status()
			if err != nil {
				return err
			}
			if enabled {
				fmt.Fprintln(cmd.OutOrStdout(), "Enabled:")
				fmt.Fprintln(cmd.OutOrStdout(), " ", line)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "Disabled.")
			}
			return nil
		},
	}
}
