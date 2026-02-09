package core

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Config holds the gvm configuration
type Config struct {
	// DownloadSource is the base URL for Go downloads (default: https://go.dev/dl/)
	DownloadSource string `json:"download_source"`
	// DownloadSourceJSON is the JSON API endpoint for version info (default: https://go.dev/dl/?mode=json&include=all)
	DownloadSourceJSON string `json:"download_source_json"`
}

const (
	// DefaultDownloadSource is the default Go download source
	DefaultDownloadSource = "https://go.dev/dl/"
	// DefaultDownloadSourceJSON is the default JSON API endpoint
	DefaultDownloadSourceJSON = "https://go.dev/dl/?mode=json&include=all"
)

// DefaultConfig returns a config with default values
func DefaultConfig() *Config {
	return &Config{
		DownloadSource:      DefaultDownloadSource,
		DownloadSourceJSON:  DefaultDownloadSourceJSON,
	}
}

// ConfigPath returns the path to the config file
func ConfigPath() (string, error) {
	d, err := GvmDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "config.json"), nil
}

// LoadConfig loads the configuration from the config file
// If the file doesn't exist, returns a default config
func LoadConfig() (*Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		// Return default config
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Ensure defaults for empty values
	if cfg.DownloadSource == "" {
		cfg.DownloadSource = DefaultDownloadSource
	}
	if cfg.DownloadSourceJSON == "" {
		cfg.DownloadSourceJSON = DefaultDownloadSourceJSON
	}

	return &cfg, nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(cfg *Config) error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0o644)
}

// GetDownloadSource returns the configured download source URL
func GetDownloadSource() (string, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return "", err
	}
	return cfg.DownloadSource, nil
}

// GetDownloadSourceJSON returns the configured JSON API endpoint
func GetDownloadSourceJSON() (string, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return "", err
	}
	return cfg.DownloadSourceJSON, nil
}
