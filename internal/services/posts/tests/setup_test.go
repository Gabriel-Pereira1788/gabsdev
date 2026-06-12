package services_test

import (
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find project root (go.mod)")
		}
		dir = parent
	}
	if err := os.Chdir(dir); err != nil {
		panic("chdir to project root: " + err.Error())
	}
}
