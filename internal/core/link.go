package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func LinkVersion(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	var goroot string
	var goBin string

	info, err := os.Stat(absPath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		// 假设用户提供的是 GOROOT 路径
		goroot = absPath
		goBin = filepath.Join(goroot, "bin", "go")
		if runtime.GOOS == "windows" {
			goBin += ".exe"
		}
		if _, err := os.Stat(goBin); os.IsNotExist(err) {
			return fmt.Errorf("invalid Go SDK path: %s not found inside %s", filepath.Base(goBin), goroot)
		}
	} else {
		// 假设用户提供的是 go 二进制文件路径
		// 尝试通过 go env GOROOT 获取真实路径
		fmt.Printf("🔍 Resolving GOROOT from binary: %s\n", absPath)
		cmd := exec.Command(absPath, "env", "GOROOT")
		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get GOROOT from %s: %v", absPath, err)
		}
		goroot = strings.TrimSpace(string(out))
		goBin = absPath
		fmt.Printf("✅ Found GOROOT: %s\n", goroot)
	}

	// 2. 获取版本号
	cmd := exec.Command(goBin, "version")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run go version: %v", err)
	}
	// Output format: "go version go1.21.5 darwin/arm64"
	parts := strings.Fields(string(out))
	if len(parts) < 3 {
		return fmt.Errorf("unknown go version output: %s", string(out))
	}
	versionStr := parts[2] // "go1.21.5"
	if strings.HasPrefix(versionStr, "go") {
		versionStr = strings.TrimPrefix(versionStr, "go")
	}

	fmt.Printf("🔍 Detected Go version: %s\n", versionStr)

	d, err := GvmDir()
	if err != nil {
		return err
	}

	// 3. 创建软链接
	// 目标: ~/.gvm/go1.21.5 -> /usr/local/go (GOROOT)
	linkName := filepath.Join(d, "go"+versionStr)

	// Check if exists
	if info, err := os.Lstat(linkName); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			// It's a symlink, check where it points
			target, _ := os.Readlink(linkName)
			if target == goroot {
				fmt.Printf("⚠️  Version %s is already linked to %s\n", versionStr, goroot)
				return nil
			}
			fmt.Printf("⚠️  Updating existing link for %s\n", versionStr)
			os.Remove(linkName)
		} else {
			return fmt.Errorf("version %s already exists and is not a symlink (it might be a real installation)", versionStr)
		}
	}

	if err := os.Symlink(goroot, linkName); err != nil {
		return fmt.Errorf("failed to create symlink: %v", err)
	}

	fmt.Printf("🔗 Linked %s -> %s\n", linkName, goroot)
	fmt.Printf("🎉 You can now use it with: gvm use %s\n", versionStr)
	return nil
}
