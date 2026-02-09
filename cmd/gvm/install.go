package gvm

import (
	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install a Go version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.InstallVersion(args[0])
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
