package core

import (
	"archive/tar"
	"compress/gzip"
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
	if strings.HasPrefix(version, "go") {
		version = strings.TrimPrefix(version, "go")
	}
	vdir := filepath.Join(d, "go"+version)
	if _, err := os.Stat(vdir); err == nil {
		return nil
	}
	osys := runtime.GOOS
	arch := runtime.GOARCH
	url := fmt.Sprintf("https://go.dev/dl/go%s.%s-%s.tar.gz", version, osys, arch)
	tarPath := filepath.Join(d, fmt.Sprintf("go%s.%s-%s.tar.gz", version, osys, arch))
	if err := os.MkdirAll(d, 0o755); err != nil {
		return err
	}
	fmt.Println("🖕Install go", "go"+version)
	fmt.Println("🌿Install from `" + url + "`")
	fmt.Println("🚀Save to:", tarPath)
	tmpf, err := os.OpenFile(tarPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("下载失败: %s", resp.Status)
	}
	cl := resp.ContentLength
	start := time.Now()
	var written int64
	buf := make([]byte, 32*1024)
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := tmpf.Write(buf[:n]); werr != nil {
				tmpf.Close()
				return werr
			}
			written += int64(n)
			printProgress(written, cl, start)
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			tmpf.Close()
			return rerr
		}
	}
	fmt.Println()
	if err := tmpf.Close(); err != nil {
		return err
	}
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
		return fmt.Errorf("包结构错误")
	}
	if err := os.Rename(src, vdir); err != nil {
		return err
	}
	_ = os.Remove(tarPath)
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
