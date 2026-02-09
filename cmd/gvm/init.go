package gvm

import (
	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gvm environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.InitEnv()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
