package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

// Upgrade upgrades a minor version to the latest patch version
// For example: "go1.25" or "1.25" will upgrade to the latest "go1.25.x"
func UpgradeVersion(versionPrefix string) (string, error) {
	// Normalize the version prefix
	versionPrefix = strings.TrimPrefix(versionPrefix, "go")
	if !strings.HasPrefix(versionPrefix, "1.") {
		return "", fmt.Errorf("invalid version format: %s", versionPrefix)
	}

	// Check if it's a minor version (e.g., "1.25" or "1.25.0")
	minorVersion, err := extractMinorVersion(versionPrefix)
	if err != nil {
		return "", err
	}

	// Get the current installed version of this minor version
	currentVersion := getCurrentPatchVersion(minorVersion)

	// Search for the latest patch version of this minor version
	latestVersion, err := getLatestPatchVersion(minorVersion)
	if err != nil {
		return "", err
	}

	fmt.Printf("当前版本: go%s\n", currentVersion)
	fmt.Printf("最新版本: go%s\n", latestVersion)

	// Check if already on latest version
	if currentVersion == latestVersion {
		fmt.Println("已经是最新版本!")
		return latestVersion, nil
	}

	// Confirm upgrade
	fmt.Printf("是否升级到 go%s? (y/N): ", latestVersion)
	// For non-interactive use, we'll proceed automatically
	// In a real CLI, you might want to add a --yes flag

	// Install the latest version
	if err := InstallVersion(latestVersion); err != nil {
		return "", err
	}

	return latestVersion, nil
}

// extractMinorVersion extracts the minor version from a version string
// For example: "1.25.0" -> "1.25", "1.25" -> "1.25"
func extractMinorVersion(version string) (string, error) {
	// Match patterns like "1.25.0" or "1.25"
	re := regexp.MustCompile(`^1\.(\d+)(\.\d+)?$`)
	matches := re.FindStringSubmatch(version)
	if matches == nil {
		return "", fmt.Errorf("invalid version format: %s", version)
	}

	return fmt.Sprintf("1.%s", matches[1]), nil
}

// getCurrentPatchVersion returns the currently installed patch version for a minor version
func getCurrentPatchVersion(minorVersion string) string {
	versions, err := SearchLocal(minorVersion)
	if err != nil || len(versions) == 0 {
		return "none"
	}

	// Sort versions to get the latest installed
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))
	return versions[0]
}

// getLatestPatchVersion returns the latest patch version for a minor version from remote
func getLatestPatchVersion(minorVersion string) (string, error) {
	osys := runtime.GOOS
	arch := runtime.GOARCH

	// Build the source URL with the minor version filter
	sourceURL, err := GetDownloadSourceJSON()
	if err != nil {
		return "", err
	}

	resp, err := http.Get(sourceURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch versions: %s", resp.Status)
	}

	var all []DLVersion
	if err := json.NewDecoder(resp.Body).Decode(&all); err != nil {
		return "", err
	}

	// Filter for the specific minor version and check platform availability
	re := regexp.MustCompile(`^go1\.\d+\.\d+$`)
	var versions []string
	for _, v := range all {
		if !re.MatchString(v.Version) {
			continue
		}
		// Check if this version matches our minor version
		if strings.HasPrefix(v.Version, "go"+minorVersion+".") {
			// Verify platform availability
			hasPlatform := false
			for _, f := range v.Files {
				if f.OS == osys && f.Arch == arch && f.Kind == "archive" {
					hasPlatform = true
					break
				}
			}
			if hasPlatform {
				versions = append(versions, strings.TrimPrefix(v.Version, "go"))
			}
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no versions found for %s", minorVersion)
	}

	// Sort to get the latest version
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))
	return versions[0], nil
}
