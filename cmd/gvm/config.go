package gvm

import (
	"encoding/json"
	"fmt"

	"github.com/ibreez3/gvm/internal/core"
	"github.com/spf13/cobra"
)

var (
	configSource      string
	configSourceJSON  string
	configShow        bool
	configReset       bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "管理 gvm 配置",
	Long: `管理 gvm 配置，包括设置 Go 下载源。

配置文件位置: ~/.gvm/config.json

可用配置项:
  download_source      Go 版本下载源 (默认: https://go.dev/dl/)
  download_source_json Go 版本 JSON API (默认: https://go.dev/dl/?mode=json&include=all)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleConfigCommand()
	},
}

func init() {
	configCmd.Flags().StringVar(&configSource, "source", "", "设置 Go 下载源 URL")
	configCmd.Flags().StringVar(&configSourceJSON, "json-source", "", "设置 Go JSON API URL")
	configCmd.Flags().BoolVar(&configShow, "show", false, "显示当前配置")
	configCmd.Flags().BoolVar(&configReset, "reset", false, "重置为默认配置")
	rootCmd.AddCommand(configCmd)
}

func handleConfigCommand() error {
	// Load current config
	cfg, err := core.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// Handle reset flag
	if configReset {
		cfg = core.DefaultConfig()
		if err := core.SaveConfig(cfg); err != nil {
			return fmt.Errorf("重置配置失败: %w", err)
		}
		fmt.Println("配置已重置为默认值")
		printConfig(cfg)
		return nil
	}

	// Handle show flag
	if configShow {
		printConfig(cfg)
		return nil
	}

	// Handle setting values
	modified := false

	if configSource != "" {
		cfg.DownloadSource = configSource
		modified = true
		fmt.Printf("设置 download_source = %s\n", configSource)
	}

	if configSourceJSON != "" {
		cfg.DownloadSourceJSON = configSourceJSON
		modified = true
		fmt.Printf("设置 download_source_json = %s\n", configSourceJSON)
	}

	// If no flags were provided, show current config
	if !modified {
		printConfig(cfg)
		return nil
	}

	// Save the modified config
	if err := core.SaveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	fmt.Println("配置已保存")
	return nil
}

func printConfig(cfg *core.Config) {
	fmt.Println("当前配置:")
	fmt.Println("================")
	b, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(b))
	fmt.Println("================")
}
