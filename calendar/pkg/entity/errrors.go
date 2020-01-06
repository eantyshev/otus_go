package appointment

import "errors"

var (
	ErrInvalidTime = errors.New("Invalid time range")
	ErrIdNotFound  = errors.New("ID is not found")
	ErrTimeBusy    = errors.New("Another entity on this period")
)
