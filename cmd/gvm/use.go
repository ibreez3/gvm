package main

import (
	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "Switch to a Go version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.UseVersion(args[0])
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
