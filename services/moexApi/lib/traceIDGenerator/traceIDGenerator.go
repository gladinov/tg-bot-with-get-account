package traceidgenerator

import (
	"main/lib/e"

	"github.com/google/uuid"
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
