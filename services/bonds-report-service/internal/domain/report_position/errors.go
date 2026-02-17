package report

import "errors"

var (
	ErrZeroQuantity = errors.New("quantity could't be zero")
	ErrUnknownOpp   = errors.New("cannot apply operation to position. Operation type is unknown")
)
