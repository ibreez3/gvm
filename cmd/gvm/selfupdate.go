package gvm

import (
	"fmt"

	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update gvm to the latest version",
	Long:  `Check and update gvm to the latest version. If already up-to-date, no action will be taken.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		checkOnly, _ := cmd.Flags().GetBool("check")

		if checkOnly {
			hasUpdate, latest, err := core.CheckUpdate()
			if err != nil {
				return err
			}
			fmt.Printf("Current version: %s\n", core.GvmVersion)
			fmt.Printf("Latest version: %s\n", latest)
			if hasUpdate {
				fmt.Println("A new version is available!")
			} else {
				fmt.Println("Already up-to-date")
			}
			return nil
		}

		return core.SelfUpdate()
	},
}

func init() {
	selfUpdateCmd.Flags().Bool("check", false, "Only check for updates, do not install")
	rootCmd.AddCommand(selfUpdateCmd)
}
