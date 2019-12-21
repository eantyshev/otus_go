package appointment

import "errors"

var (
	ErrInvalidTime = errors.New("Invalid time range")
	ErrConflictId  = errors.New("ID already exists")
	ErrIdNotFound  = errors.New("ID is not found")
	ErrTimeBusy    = errors.New("Another entity on this period")
)
