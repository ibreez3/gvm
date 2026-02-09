package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []Asset `json:"assets"`
}

// Asset represents a release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

const (
	// GvmGitHubRepo is the GitHub repository for gvm
	GvmGitHubRepo = "ibreez3/gvm"
	// GitHubAPIURL is the GitHub API URL for releases
	GitHubAPIURL = "https://api.github.com/repos/" + GvmGitHubRepo + "/releases/latest"
)

// LatestVersion fetches the latest version of gvm from GitHub releases
func LatestVersion() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", GitHubAPIURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch latest version: %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

// GvmVersion returns the current version of gvm
// This is set during build using ldflags
var GvmVersion = "dev"

// SelfUpdate updates gvm to the latest version
func SelfUpdate() error {
	fmt.Printf("Current version: %s\n", GvmVersion)

	// Get latest version from GitHub
	latest, err := LatestVersion()
	if err != nil {
		return err
	}
	fmt.Printf("Latest version: %s\n", latest)

	// Check if already up to date
	if GvmVersion != "dev" && GvmVersion == latest {
		fmt.Println("Already up to date!")
		return nil
	}

	// Get the binary path
	binPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Determine the asset name based on OS and architecture
	assetName, err := getAssetName()
	if err != nil {
		return err
	}

	fmt.Printf("Downloading %s...\n", assetName)

	// Get the download URL for the asset
	downloadURL, err := getAssetDownloadURL(latest, assetName)
	if err != nil {
		return err
	}

	// Download to a temporary file
	tmpDir := os.TempDir()
	tmpPath := filepath.Join(tmpDir, "gvm-new-"+assetName)

	fmt.Printf("Downloading from %s...\n", downloadURL)
	if err := downloadUpdateFile(downloadURL, tmpPath); err != nil {
		return err
	}
	defer os.Remove(tmpPath)

	fmt.Println("Download complete!")

	// Replace the old binary with the new one
	fmt.Printf("Updating %s...\n", binPath)

	// On Unix systems, we need to remove the old file first
	// On Windows, we need to move the old file and then replace it
	if err := replaceBinary(tmpPath, binPath); err != nil {
		return err
	}

	fmt.Printf("Successfully updated to %s!\n", latest)
	return nil
}

func getAssetName() (string, error) {
	osys := runtime.GOOS
	arch := runtime.GOARCH

	// Map arch to the naming convention used in releases
	var archStr string
	switch arch {
	case "amd64":
		archStr = "x86_64"
	case "arm64":
		archStr = "arm64"
	case "386":
		archStr = "i386"
	default:
		return "", fmt.Errorf("unsupported architecture: %s", arch)
	}

	// Determine the extension based on OS
	var ext string
	switch osys {
	case "windows":
		ext = ".exe"
	default:
		ext = ""
	}

	return fmt.Sprintf("gvm_%s_%s%s", osys, archStr, ext), nil
}

func getAssetDownloadURL(version, assetName string) (string, error) {
	// Get release info from GitHub
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", GvmGitHubRepo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch release info: %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	// Find the matching asset
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("asset not found: %s", assetName)
}

func downloadUpdateFile(url, dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	// Show progress
	size := resp.ContentLength
	var written int64
	buf := make([]byte, 32*1024)
	start := time.Now()

	for {
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			nw, ew := out.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				return ew
			}
			if nr != nw {
				return io.ErrShortWrite
			}
			printUpdateProgress(written, size, start)
		}
		if er != nil {
			if er != io.EOF {
				return er
			}
			break
		}
	}
	fmt.Println()
	return nil
}

func printUpdateProgress(written, total int64, start time.Time) {
	pct := float64(0)
	if total > 0 {
		pct = float64(written) / float64(total)
	}
	speed := float64(written) / time.Since(start).Seconds()
	eta := "--"
	if speed > 0 && total > 0 {
		rem := float64(total-written) / speed
		eta = formatDuration(time.Duration(rem) * time.Second)
	}
	left := formatSize(written)
	right := formatSize(total)
	bar := progressBar(pct, 30)
	fmt.Printf("\r  %s / %s %s  %5.2f%% %s/s %s", left, right, bar, pct*100, formatSize(int64(speed)), eta)
}

func replaceBinary(src, dest string) error {
	// Make the new file executable
	if err := os.Chmod(src, 0o755); err != nil {
		return err
	}

	// On Windows, we need to handle the file replacement differently
	if strings.HasPrefix(runtime.GOOS, "windows") {
		// Move the old file to a temporary location
		oldPath := dest + ".old"
		if err := os.Rename(dest, oldPath); err != nil {
			return err
		}
		// Move the new file to the destination
		if err := os.Rename(src, dest); err != nil {
			// Try to restore the old file
			_ = os.Rename(oldPath, dest)
			return err
		}
		// Clean up the old file
		_ = os.Remove(oldPath)
	} else {
		// On Unix systems, we can just rename
		if err := os.Rename(src, dest); err != nil {
			return err
		}
	}

	return nil
}

// CheckUpdate checks if there's a newer version available
func CheckUpdate() (bool, string, error) {
	latest, err := LatestVersion()
	if err != nil {
		return false, "", err
	}

	if GvmVersion == "dev" {
		return true, latest, nil
	}

	if GvmVersion != latest {
		return true, latest, nil
	}

	return false, latest, nil
}
