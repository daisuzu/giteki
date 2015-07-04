package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	downloadScript = "src/github.com/daisuzu/giteki/scripts/download.py"
	readScript     = "src/github.com/daisuzu/giteki/scripts/read.py"
)

func getScript(s string) (string, error) {
	for _, path := range filepath.SplitList(os.Getenv("GOPATH")) {
		script := filepath.Join(path, s)
		if _, err := os.Stat(script); err == nil {
			return script, nil
		}
	}
	return "", fmt.Errorf("Not found `%s`", s)
}

func GetDownloadScript() (string, error) {
	return getScript(downloadScript)
}

func GetReadScript() (string, error) {
	return getScript(readScript)
}
