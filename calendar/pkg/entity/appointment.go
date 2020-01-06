package appointment

import (
	"github.com/google/uuid"
	"time"
)

type Appointment struct {
	Uuid        uuid.UUID `json:"uuid" validate:"required"`
	Summary     string    `json:"summary" validate:"required"`
	Description string    `json:"description"`
	TimeStart   time.Time `json:"time_start" validate:"required"`
	TimeEnd     time.Time `json:"time_end" validate:"required"`
	Owner       string    `json:"owner" validate:"required"`
}
