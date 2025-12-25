package traceidgenerator

import (
	"github.com/google/uuid"
	"main.go/lib/e"
)

func New() (string, error) {
	const op = "traceidgenerator.New"
	id, err := uuid.NewRandom()
	if err != nil {
		return "", e.WrapIfErr(op, err)
	}
	return id.String(), nil
}

func Must() string {
	id, err := New()
	if err != nil {
		panic(err)
	}
	return id
}
