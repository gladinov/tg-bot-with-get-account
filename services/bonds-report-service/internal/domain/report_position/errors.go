package report

import "errors"

var (
	ErrZeroQuantity = errors.New("quantity could't be zero")
	ErrUnknownOpp   = errors.New("cannot apply operation to position. Operation type is unknown")
	ErrZeroDivision = errors.New("division could't be zero")
	ErrInvalidDate  = errors.New("buy date could't be after sell date")
)
