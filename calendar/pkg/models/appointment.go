package models

import "time"

type Appointment struct {
	ID              int64          `json:"id" validate:"required"`
	Summary         string         `json:"summary" validate:"required"`
	Description     string         `json:"description"`
	StartsAt        time.Time      `json:"starts_at"`
	DurationMinutes uint16         `json:"duration_minutes,omitempty"`
	IsRegular       bool           `json:"is_regular,omitempty"`
	DaysOfWeek      []time.Weekday `json:"days_of_week,omitempty"`
}
