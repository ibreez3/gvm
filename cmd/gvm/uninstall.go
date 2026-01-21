package main

import (
	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [version]",
	Short: "Uninstall a Go version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.UninstallVersion(args[0])
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
