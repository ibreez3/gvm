package main

import (
	"fmt"
	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed Go versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		remote, _ := cmd.Flags().GetBool("remote")
		if remote {
			versions, err := core.ListRemote(20)
			if err != nil {
				return err
			}
			for _, v := range versions {
				fmt.Println(v)
			}
			return nil
		}
		
		versions, err := core.ListLocal()
		if err != nil {
			return err
		}
		current, _ := core.CurrentVersion()
		for _, v := range versions {
			if v == current {
				fmt.Printf("* %s\n", v)
			} else {
				fmt.Printf("  %s\n", v)
			}
		}
		return nil
	},
}

func init() {
	listCmd.Flags().BoolP("remote", "r", false, "List remote versions")
	rootCmd.AddCommand(listCmd)
}
