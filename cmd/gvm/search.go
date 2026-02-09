package gvm

import (
	"fmt"
	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Search for Go versions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		versions, err := core.SearchRemote(args[0], 20)
		if err != nil {
			return err
		}
		for _, v := range versions {
			fmt.Println(v)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
