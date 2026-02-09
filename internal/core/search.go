package core

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "regexp"
    "strings"
)

func SearchRemote(prefix string, limit int) ([]string, error) {
    prefix = strings.TrimPrefix(prefix, "go")
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
        return nil, fmt.Errorf("远程查询失败: %s", resp.Status)
    }
    var all []DLVersion
    dec := json.NewDecoder(resp.Body)
    if err := dec.Decode(&all); err != nil {
        return nil, err
    }
    re := regexp.MustCompile(`^go\d+\.\d+\.\d+$`)
    var vv []string
    for _, it := range all {
        if !re.MatchString(it.Version) {
            continue
        }
        v := strings.TrimPrefix(it.Version, "go")
        if strings.HasPrefix(v, prefix+".") {
            vv = append(vv, v)
            if limit > 0 && len(vv) >= limit {
                break
            }
        }
    }
    return vv, nil
}

func SearchLocal(prefix string) ([]string, error) {
    prefix = strings.TrimPrefix(prefix, "go")
    d, err := GvmDir()
    if err != nil {
        return nil, err
    }
    es, err := os.ReadDir(d)
    if err != nil {
        return nil, err
    }
    re := regexp.MustCompile(`^go\d+\.\d+\.\d+$`)
    var vv []string
    for _, e := range es {
        if e.IsDir() && re.MatchString(e.Name()) {
            v := strings.TrimPrefix(e.Name(), "go")
            if strings.HasPrefix(v, prefix+".") {
                vv = append(vv, v)
            }
        }
    }
    return vv, nil
}

