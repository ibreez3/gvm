package core

import (
    "os"
    "path/filepath"
    "strings"
)

func UseVersion(version string) error {
    d, err := GvmDir()
    if err != nil {
        return err
    }
    if strings.HasPrefix(version, "go") {
        version = strings.TrimPrefix(version, "go")
    }
    vdir := filepath.Join(d, "go"+version)
    if _, err := os.Stat(vdir); err != nil {
        return err
    }
    link := filepath.Join(d, "goroot")
    if fi, err := os.Lstat(link); err == nil && fi.Mode()&os.ModeSymlink != 0 {
        _ = os.Remove(link)
    } else if err == nil {
        _ = os.Remove(link)
    }
    if err := os.Symlink(vdir, link); err != nil {
        return err
    }
    return nil
}

func CurrentVersion() (string, error) {
    d, err := GvmDir()
    if err != nil {
        return "", err
    }
    link := filepath.Join(d, "goroot")
    t, err := os.Readlink(link)
    if err != nil {
        return "", err
    }
    b := filepath.Base(t)
    if strings.HasPrefix(b, "go") {
        return strings.TrimPrefix(b, "go"), nil
    }
    return b, nil
}

