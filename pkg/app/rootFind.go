package rootfind

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var Root string

func MustInitialize() {
	var err error
	Root, err = findProjectRoot()
	if err != nil {
		panic(fmt.Sprintf("Cannot initialize application: %v", err))
	}
}

func MustGetRoot() string {
	if Root == "" {
		panic("Application not initialized. Call MustInitialize() first")
	}
	return Root
}

func findProjectRoot() (string, error) {
	if root := findRootBySource(); root != "" {
		return root, nil
	}

	if root := findRootByExecutable(); root != "" {
		return root, nil
	}

	return "", fmt.Errorf("cannot locate project root (go.mod not found)")
}

func findRootBySource() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	currentDir := filepath.Dir(filename)
	for {
		if isProjectRoot(currentDir) {
			return currentDir
		}
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}
	return ""
}

func findRootByExecutable() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}

	currentDir := filepath.Dir(exe)
	for i := 0; i < 10; i++ {
		if isProjectRoot(currentDir) {
			return currentDir
		}
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}
	return ""
}

func isProjectRoot(dir string) bool {
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		return true
	}
	return false
}
