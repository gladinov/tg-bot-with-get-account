package domain

import "errors"

var (
	ErrEmptyUids     = errors.New("no uids")
	ErrNoCurrency    = errors.New("no currency")
	ErrNoOpperations = errors.New("no operations")
	ErrEmptyReport   = errors.New("no elements in report")
)

var ErrEmptyUidAfterUpdate = errors.New("asset uid by this instrument uid is not exist")
