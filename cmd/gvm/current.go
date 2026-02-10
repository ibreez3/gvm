package gvm

import (
	"fmt"

	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current version",
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := core.CurrentVersion()
		if err != nil {
			return err
		}
		fmt.Println(v)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
