package core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func HomeDir() (string, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return h, nil
}

func GvmDir() (string, error) {
	h, err := HomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(h, ".gvm"), nil
}

func InitEnv() error {
	d, err := GvmDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(d, 0o755); err != nil {
		return err
	}
	f := filepath.Join(d, ".gvmrc")
	content := strings.Join([]string{
		"export GOROOT=$HOME/.gvm/goroot",
		"export PATH=$PATH:$GOROOT/bin",
		"export GOPATH=$HOME/go",
		"export GOBIN=$GOPATH/bin",
		"export PATH=$PATH:$GGOBIN",
		"export GOPROXY=https://goproxy.cn,direct",
		"",
	}, "\n")
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		return err
	}
	rc, err := detectShellRC()
	if err != nil {
		return err
	}
	if rc != "" {
		line := "if [ -f \"$HOME/.gvm/.gvmrc\" ]; then\n    source \"$HOME/.gvm/.gvmrc\"\nfi\n"
		b, _ := os.ReadFile(rc)
		s := string(b)
		if !strings.Contains(s, "/.gvm/.gvmrc") {
			f2, err := os.OpenFile(rc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				return err
			}
			defer f2.Close()
			if _, err := f2.WriteString("\n# gvm shell setup\n" + line); err != nil {
				return err
			}
		}
	}
	c := filepath.Join(d, "config.toml")
	if _, err := os.Stat(c); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(c, []byte(""), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func detectShellRC() (string, error) {
	h, err := HomeDir()
	if err != nil {
		return "", err
	}
	z := filepath.Join(h, ".zshrc")
	b := filepath.Join(h, ".bashrc")
	if _, err := os.Stat(z); err == nil {
		return z, nil
	}
	if _, err := os.Stat(b); err == nil {
		return b, nil
	}
	return "", nil
}
