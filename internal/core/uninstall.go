package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// UninstallBatchSpec defines the batch uninstall specification
type UninstallBatchSpec struct {
	Below     string // Uninstall versions below this (exclusive)
	Pattern   string // Uninstall versions matching this pattern
	Keep      int    // Keep this many latest versions
	All       bool   // Uninstall all versions
	KeepCurrent bool // Keep the currently active version
}

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

// UninstallBatch performs batch uninstall based on the specification
func UninstallBatch(spec *UninstallBatchSpec) ([]string, error) {
	versions, err := ListLocal()
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no versions installed")
	}

	// Get current version to protect it if needed
	currentVersion := ""
	if spec.KeepCurrent {
		currentVersion, _ = CurrentVersion()
	}

	// Filter versions based on spec
	var toUninstall []string
	for _, v := range versions {
		// Skip current version if requested
		if spec.KeepCurrent && v == currentVersion {
			continue
		}

		// Apply filters
		shouldUninstall := false

		if spec.All {
			shouldUninstall = true
		} else if spec.Below != "" {
			// Compare versions
			if compareVersions(v, spec.Below) < 0 {
				shouldUninstall = true
			}
		} else if spec.Pattern != "" {
			// Match pattern (supports wildcards like "1.21.*")
			if matchVersionPattern(v, spec.Pattern) {
				shouldUninstall = true
			}
		} else if spec.Keep > 0 {
			// This is handled differently - we need to keep N latest versions
			// We'll sort and keep the top N
		}

		if shouldUninstall {
			toUninstall = append(toUninstall, v)
		}
	}

	// Handle "keep N latest" logic
	if spec.Keep > 0 && spec.Below == "" && spec.Pattern == "" && !spec.All {
		// Sort versions (newest first)
		sortedVersions := make([]string, len(versions))
		copy(sortedVersions, versions)
		sort.Slice(sortedVersions, func(i, j int) bool {
			return compareVersions(sortedVersions[i], sortedVersions[j]) > 0
		})

		// Versions to uninstall are those beyond the first N
		for i := spec.Keep; i < len(sortedVersions); i++ {
			v := sortedVersions[i]
			// Skip current version if requested
			if spec.KeepCurrent && v == currentVersion {
				continue
			}
			toUninstall = append(toUninstall, v)
		}
	}

	if len(toUninstall) == 0 {
		return nil, fmt.Errorf("no versions match the criteria")
	}

	// Show what will be uninstalled
	fmt.Printf("将卸载以下版本:\n")
	for _, v := range toUninstall {
		fmt.Printf("  - go%s\n", v)
	}

	// Perform uninstall
	var uninstalled []string
	for _, v := range toUninstall {
		if err := UninstallVersion(v); err != nil {
			fmt.Printf("⚠️  Failed to uninstall go%s: %v\n", v, err)
		} else {
			uninstalled = append(uninstalled, v)
		}
	}

	return uninstalled, nil
}

// compareVersions compares two version strings (e.g., "1.22.0" vs "1.21.0")
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
func compareVersions(v1, v2 string) int {
	// Normalize versions (remove "go" prefix if present)
	v1 = strings.TrimPrefix(v1, "go")
	v2 = strings.TrimPrefix(v2, "go")

	// Split by "."
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Compare each part
	for i := 0; i < 3; i++ {
		// Default to 0 if missing part
		var n1, n2 int
		if i < len(parts1) {
			n1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			n2, _ = strconv.Atoi(parts2[i])
		}

		if n1 > n2 {
			return 1
		}
		if n1 < n2 {
			return -1
		}
	}

	return 0
}

// matchVersionPattern checks if a version matches a pattern
// Supports wildcards: "1.21.*" matches "1.21.0", "1.21.1", etc.
func matchVersionPattern(version, pattern string) bool {
	// Normalize
	version = strings.TrimPrefix(version, "go")
	pattern = strings.TrimPrefix(pattern, "go")

	// Handle wildcard
	if strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, ".*")
		return strings.HasPrefix(version, prefix+".")
	}

	// Exact match
	return version == pattern
}
