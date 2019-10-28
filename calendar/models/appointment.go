package models

import "time"

type WeekDay byte

const (
	Sunday WeekDay = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type Appointment struct {
	ID              int64     `json:"id"`
	Summary         string    `json:"summary" validate:"required"`
	Description     string    `json:"description"`
	StartsAt        time.Time `json:"starts_at"`
	DurationMinutes uint16    `json:"duration_minutes"`
	IsRegular       bool      `json:"is_regular"`
	DaysOfWeek      []WeekDay `json:"days_of_week"`
}
