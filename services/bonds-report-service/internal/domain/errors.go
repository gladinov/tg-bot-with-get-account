package domain

import "errors"

var (
	ErrEmptyUids     = errors.New("no uids")
	ErrNoCurrency    = errors.New("no currency")
	ErrNoOpperations = errors.New("no operations")
	ErrEmptyReport   = errors.New("no elements in report")
)

var ErrEmptyUidAfterUpdate = errors.New("asset uid by this instrument uid is not exist")

var (
	ErrCloseAccount            = errors.New("close account haven't portffolio positions")
	ErrNoAcces                 = errors.New("this token no access to account")
	ErrEmptyAccountIdInRequest = errors.New("accountId could not be empty")
	ErrUnspecifiedAccount      = errors.New("account is unspecified")
	ErrNewNotOpenYetAccount    = errors.New("accountId is not opened yet")
	ErrEmptyInstrumentUid      = errors.New("instrumentUid could not be empty string")
	ErrEmptyFigi               = errors.New("figi could not be empty string")
	ErrEmptyQuery              = errors.New("query could not be empty")
	ErrEmptyUid                = errors.New("uid could not be empty string")
	ErrEmptyPositionUid        = errors.New("positionUid could not be empty string")
)

var (
	ErrEmptyAccountID               = errors.New("accountID could not be empty string")
	ErrInvalidFromDate              = errors.New("from can't be more than the current date")
	ErrEmptyInstrumentShortResponce = errors.New("instrument short responce is empty")
	ErrInstrumentNotShare           = errors.New("instrument is not share")
	ErrEmptyTicker                  = errors.New("ticker is empty")
)
