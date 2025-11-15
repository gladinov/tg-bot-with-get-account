package pathwd

import (
	"os"
	"path/filepath"

	"main.go/lib/e"
)

func PathFromWD(workDir string, path string) (string, error) {
	pathWD := filepath.Join(workDir, path)
	err := os.MkdirAll(filepath.Dir(pathWD), 0755)
	if err != nil {
		return "", e.WrapIfErr("can't create directory", err)
	}
	return pathWD, nil
}
