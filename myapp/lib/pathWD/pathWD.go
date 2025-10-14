package pathwd

import (
	"os"
	"path/filepath"

	"main.go/lib/e"
)

func PathFromWD(path string) (string, error) {
	// TODO: Обрабтать ошибку с точками (./) в путях
	cwd, _ := os.Getwd()
	pathWD := filepath.Join(cwd, path)
	err := os.MkdirAll(filepath.Dir(pathWD), 0755)
	if err != nil {
		return "", e.WrapIfErr("can't create directory", err)
	}
	return pathWD, nil
}
