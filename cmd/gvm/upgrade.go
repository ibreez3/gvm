package gvm

import (
	"fmt"

	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var (
	upgradeUse   bool
	upgradeYes   bool
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [version]",
	Short: "升级 Go 次版本到最新的补丁版本",
	Long: `升级指定次版本到最新的补丁版本。

示例:
  gvm upgrade 1.25    # 升级到最新的 1.25.x 版本
  gvm upgrade go1.25  # 同上 (支持 go 前缀)`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := core.UpgradeVersion(args[0])
		if err != nil {
			return err
		}

		// Automatically use the newly upgraded version if requested
		if upgradeUse {
			if err := core.UseVersion(version); err != nil {
				return err
			}
			fmt.Printf("已切换到 go%s\n", version)
		}

		return nil
	},
}

func init() {
	upgradeCmd.Flags().BoolVarP(&upgradeUse, "use", "u", false, "升级后自动切换到新版本")
	upgradeCmd.Flags().BoolVarP(&upgradeYes, "yes", "y", false, "自动确认升级")
	rootCmd.AddCommand(upgradeCmd)
}
