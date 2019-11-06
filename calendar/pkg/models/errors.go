package models

import "errors"

var (
	ErrConflictId = errors.New("ID already exists")
	ErrIdNotFound = errors.New("ID is not found")
	ErrTimeBusy   = errors.New("Another appointment on this period")
)
