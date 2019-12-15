package repository

import (
	"github.com/eantyshev/otus_go/calendar/pkg/appointment"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	Fetch(timeBegin time.Time, num int) ([]*appointment.Appointment, time.Time, error)
	GetById(id uuid.UUID) (*appointment.Appointment, error)
	Update(ap *appointment.Appointment) error
	Store(ap *appointment.Appointment) error
	Delete(id uuid.UUID) error
}
