package main

import (
	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link [path]",
	Short: "Link an external Go SDK to gvm",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.LinkVersion(args[0])
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
