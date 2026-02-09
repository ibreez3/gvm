package gvm

import (
	"fmt"

	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var (
	uninstallBelow     string
	uninstallPattern   string
	uninstallKeep      int
	uninstallAll       bool
	uninstallKeepCurrent bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [version]",
	Short: "卸载 Go 版本",
	Long: `卸载 Go 版本，支持单个或批量卸载。

示例:
  gvm uninstall 1.21.0           # 卸载单个版本
  gvm uninstall --below 1.22     # 卸载低于 1.22 的所有版本
  gvm uninstall --pattern "1.21.*"  # 卸载 1.21.x 系列的所有版本
  gvm uninstall --keep 2         # 只保留最新的 2 个版本，卸载其余
  gvm uninstall --all            # 卸载所有版本

注意: 使用批量卸载时会自动跳过当前正在使用的版本。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Batch mode: no version specified, use flags
		if len(args) == 0 {
			spec := &core.UninstallBatchSpec{
				Below:       uninstallBelow,
				Pattern:     uninstallPattern,
				Keep:        uninstallKeep,
				All:         uninstallAll,
				KeepCurrent: uninstallKeepCurrent,
			}

			// Validate that exactly one batch option is specified
			count := 0
			if uninstallBelow != "" {
				count++
			}
			if uninstallPattern != "" {
				count++
			}
			if uninstallKeep > 0 {
				count++
			}
			if uninstallAll {
				count++
			}

			if count != 1 {
				return fmt.Errorf("请指定一个批量卸载选项: --below, --pattern, --keep, 或 --all")
			}

			uninstalled, err := core.UninstallBatch(spec)
			if err != nil {
				return err
			}
			fmt.Printf("\n✅ 成功卸载 %d 个版本\n", len(uninstalled))
			return nil
		}

		// Single version mode
		if uninstallBelow != "" || uninstallPattern != "" || uninstallKeep > 0 || uninstallAll {
			return fmt.Errorf("不能同时指定版本和批量卸载选项")
		}

		return core.UninstallVersion(args[0])
	},
}

func init() {
	uninstallCmd.Flags().StringVar(&uninstallBelow, "below", "", "卸载低于此版本的所有版本")
	uninstallCmd.Flags().StringVar(&uninstallPattern, "pattern", "", "卸载匹配此模式的所有版本 (支持通配符，如 1.21.*)")
	uninstallCmd.Flags().IntVar(&uninstallKeep, "keep", 0, "只保留最新的 N 个版本，卸载其余")
	uninstallCmd.Flags().BoolVar(&uninstallAll, "all", false, "卸载所有版本")
	uninstallCmd.Flags().BoolVarP(&uninstallKeepCurrent, "keep-current", "c", true, "保留当前正在使用的版本")
	rootCmd.AddCommand(uninstallCmd)
}
