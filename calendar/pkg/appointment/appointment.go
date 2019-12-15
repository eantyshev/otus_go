package appointment

import (
	"github.com/google/uuid"
	"time"
)

type Appointment struct {
	Uuid        uuid.UUID     `json:"uuid" validate:"required"`
	Summary     string        `json:"summary" validate:"required"`
	Description string        `json:"description"`
	StartsAt    time.Time     `json:"starts_at"`
	Duration    time.Duration `json:"duration,omitempty"`
	Owner       string        `json:"owner" validate:"required"`
}
