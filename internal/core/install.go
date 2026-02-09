package core

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
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

func InstallVersion(version string) error {
	d, err := GvmDir()
	if err != nil {
		return err
	}
	version = strings.TrimPrefix(version, "go")
	vdir := filepath.Join(d, "go"+version)
	if _, err := os.Stat(vdir); err == nil {
		return fmt.Errorf("version %s already installed", version)
	}

	osys := runtime.GOOS
	arch := runtime.GOARCH

	// 1. 获取版本信息（URL 和 Checksum）
	fmt.Printf("🔍 Searching for version %s ...\n", version)
	fileInfo, err := getVersionInfo("go"+version, osys, arch)
	if err != nil {
		// Fallback: 如果 JSON 中找不到，尝试直接构造 URL（但不校验 checksum，或者给警告）
		// 为了安全，这里我们先强制要求找到，或者打印警告
		fmt.Printf("⚠️  Warning: Could not find version info in official JSON API: %v\n", err)
		fmt.Println("⚠️  Proceeding with direct download (NO CHECKSUM VERIFICATION)")
		// 构造默认 URL
		fileInfo = &File{
			Filename: fmt.Sprintf("go%s.%s-%s.tar.gz", version, osys, arch),
			SHA256:   "", // Empty means no verification
		}
		// URL 需手动构造，因为 fileInfo 只有文件名
	}

	sourceURL, err := GetDownloadSource()
	if err != nil {
		return err
	}
	// Ensure the source URL ends with a slash
	if !strings.HasSuffix(sourceURL, "/") {
		sourceURL += "/"
	}
	downloadURL := sourceURL + fileInfo.Filename
	tarPath := filepath.Join(d, fileInfo.Filename)

	if err := os.MkdirAll(d, 0o755); err != nil {
		return err
	}

	fmt.Println("⬇️  Downloading go" + version + "...")
	fmt.Println("🔗 Source:", downloadURL)
	fmt.Println("📦 Dest:", tarPath)

	// 2. 下载文件
	if err := downloadFile(downloadURL, tarPath); err != nil {
		return err
	}

	// 3. 校验 Checksum
	if fileInfo.SHA256 != "" {
		fmt.Println("🛡️  Verifying checksum...")
		if err := verifyChecksum(tarPath, fileInfo.SHA256); err != nil {
			os.Remove(tarPath) // 删除损坏的文件
			return fmt.Errorf("checksum verification failed: %v", err)
		}
		fmt.Println("✅ Checksum verified")
	} else {
		fmt.Println("⚠️  Skipping checksum verification (not available)")
	}

	// 4. 解压安装
	fmt.Println("📦 Extracting...")
	tdir, err := os.MkdirTemp("", "go-tgz-untar-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tdir)

	if err := untar(tarPath, tdir); err != nil {
		return err
	}

	src := filepath.Join(tdir, "go")
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("package structure error: 'go' directory not found")
	}

	if err := os.Rename(src, vdir); err != nil {
		return err
	}

	_ = os.Remove(tarPath)
	fmt.Printf("🎉 Successfully installed go%s\n", version)
	return nil
}

func getVersionInfo(version, osys, arch string) (*File, error) {
	// 查询包含所有版本的 JSON
	sourceURL, err := GetDownloadSourceJSON()
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(sourceURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch versions: %s", resp.Status)
	}

	var versions []DLVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, err
	}

	for _, v := range versions {
		if v.Version == version {
			for _, f := range v.Files {
				if f.OS == osys && f.Arch == arch && f.Kind == "archive" {
					return &f, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("version not found in official list")
}

func downloadFile(url, dest string) error {
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
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

	// Progress bar setup
	cl := resp.ContentLength
	start := time.Now()
	var written int64
	buf := make([]byte, 32*1024)

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
			printProgress(written, cl, start)
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

func verifyChecksum(path, expected string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	actual := hex.EncodeToString(h.Sum(nil))
	if actual != expected {
		return fmt.Errorf("expected %s, got %s", expected, actual)
	}
	return nil
}

func untar(tgz string, dest string) error {
	f, err := os.Open(tgz)
	if err != nil {
		return err
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		p := filepath.Join(dest, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(p, os.FileMode(hdr.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				return err
			}
			of, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(of, tr); err != nil {
				of.Close()
				return err
			}
			of.Close()
		case tar.TypeSymlink:
			if err := os.Symlink(hdr.Linkname, p); err != nil {
				return err
			}
		default:
		}
	}
	return nil
}

func printProgress(written int64, total int64, start time.Time) {
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
	etaPad := padRight(eta, 10)
	fmt.Printf("\r  %s / %s %s  %5.2f%% %s/s %s", left, right, bar, pct*100, formatSize(int64(speed)), etaPad)
}

func progressBar(pct float64, width int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	filled := int(pct * float64(width))
	if filled < 0 {
		filled = 0
	}
	if filled > width-1 {
		filled = width - 1
	}
	return "[" + strings.Repeat("-", filled) + ">" + strings.Repeat("-", width-filled-1) + "]"
}

func formatSize(b int64) string {
	kb := float64(b) / 1024
	if kb < 1024 {
		return fmt.Sprintf("%5.2f KB", kb)
	}
	mb := kb / 1024
	return fmt.Sprintf("%5.2f MB", mb)
}

func formatDuration(d time.Duration) string {
	s := int(d.Seconds())
	if s < 60 {
		return fmt.Sprintf("%ds", s)
	}
	m := s / 60
	s = s % 60
	if m < 60 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	h := m / 60
	m = m % 60
	return fmt.Sprintf("%dh%dm%ds", h, m, s)
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
