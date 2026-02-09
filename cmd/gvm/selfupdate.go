package gvm

import (
	"fmt"

	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "更新 gvm 到最新版本",
	Long:  `检查并更新 gvm 到最新版本。如果当前版本已经是最新，则不会进行更新。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		checkOnly, _ := cmd.Flags().GetBool("check")

		if checkOnly {
			hasUpdate, latest, err := core.CheckUpdate()
			if err != nil {
				return err
			}
			fmt.Printf("当前版本: %s\n", core.GvmVersion)
			fmt.Printf("最新版本: %s\n", latest)
			if hasUpdate {
				fmt.Println("有新版本可用!")
			} else {
				fmt.Println("已经是最新版本")
			}
			return nil
		}

		return core.SelfUpdate()
	},
}

func init() {
	selfUpdateCmd.Flags().Bool("check", false, "仅检查是否有新版本，不进行更新")
	rootCmd.AddCommand(selfUpdateCmd)
}
