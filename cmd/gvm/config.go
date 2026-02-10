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
	Short: "Manage gvm configuration",
	Long: `Manage gvm configuration, including setting Go download source.

Configuration file location: ~/.gvm/config.json

Available options:
  download_source      Go version download source (default: https://go.dev/dl/)
  download_source_json Go version JSON API (default: https://go.dev/dl/?mode=json&include=all)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleConfigCommand()
	},
}

func init() {
	configCmd.Flags().StringVar(&configSource, "source", "", "Set Go download source URL")
	configCmd.Flags().StringVar(&configSourceJSON, "json-source", "", "Set Go JSON API URL")
	configCmd.Flags().BoolVar(&configShow, "show", false, "Show current configuration")
	configCmd.Flags().BoolVar(&configReset, "reset", false, "Reset to default configuration")
	rootCmd.AddCommand(configCmd)
}

func handleConfigCommand() error {
	// Load current config
	cfg, err := core.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Handle reset flag
	if configReset {
		cfg = core.DefaultConfig()
		if err := core.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to reset config: %w", err)
		}
		fmt.Println("Configuration reset to default")
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
		fmt.Printf("Set download_source = %s\n", configSource)
	}

	if configSourceJSON != "" {
		cfg.DownloadSourceJSON = configSourceJSON
		modified = true
		fmt.Printf("Set download_source_json = %s\n", configSourceJSON)
	}

	// If no flags were provided, show current config
	if !modified {
		printConfig(cfg)
		return nil
	}

	// Save the modified config
	if err := core.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("Configuration saved")
	return nil
}

func printConfig(cfg *core.Config) {
	fmt.Println("Current configuration:")
	fmt.Println("================")
	b, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(b))
	fmt.Println("================")
}
