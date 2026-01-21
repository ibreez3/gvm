package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func UninstallVersion(version string) error {
	d, err := GvmDir()
	if err != nil {
		return err
	}
	if strings.HasPrefix(version, "go") {
		version = strings.TrimPrefix(version, "go")
	}
	vdir := filepath.Join(d, "go"+version)

	// Check if exists
	if _, err := os.Stat(vdir); os.IsNotExist(err) {
		return fmt.Errorf("version %s is not installed", version)
	}

	// Check if current
	current, err := CurrentVersion()
	if err == nil && current == version {
		fmt.Printf("⚠️  Warning: Version %s is currently in use.\n", version)
		fmt.Println("   Unlinking current version...")
		link := filepath.Join(d, "goroot")
		_ = os.Remove(link)
	}

	fmt.Printf("🗑️  Uninstalling go%s...\n", version)
	if err := os.RemoveAll(vdir); err != nil {
		return fmt.Errorf("failed to uninstall: %v", err)
	}

	fmt.Printf("✅ Successfully uninstalled go%s\n", version)
	return nil
}
