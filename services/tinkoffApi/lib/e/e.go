package e

import (
	"fmt"
	"strings"
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}

	return Wrap(msg, err)
}

type ErrResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func SplitErr(err error) (ErrResponse, error) {
	errStr := err.Error()
	components := strings.Split(errStr, "=")
	if len(components) == 1 {
		return ErrResponse{}, fmt.Errorf("haven't symbol(=) to split")
	}
	var errResp ErrResponse
	errResp.Status = strings.TrimSpace(components[len(components)-1])
	errResp.Message = strings.TrimSpace(components[len(components)-2])

	return errResp, nil

}

func IsTimeError(inputErr error) bool {
	if inputErr == nil {
		return false
	}

	errorStr := inputErr.Error()
	return strings.Contains(errorStr, "30070")
}
