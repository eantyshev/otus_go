package interfaces

import (
	"context"
	"github.com/eantyshev/otus_go/calendar/pkg/entity"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	ListOwnerPeriod(ctx context.Context, owner string, timeFrom, timeTo time.Time) ([]*appointment.Appointment, error)
	GetById(ctx context.Context, id *uuid.UUID) (*appointment.Appointment, error)
	Update(ctx context.Context, ap *appointment.Appointment) error
	Create(ctx context.Context, ap *appointment.Appointment) (*uuid.UUID, error)
	Delete(ctx context.Context, id *uuid.UUID) error
	FetchPeriod(ctx context.Context, timeFrom, timeTo time.Time) ([]*appointment.Appointment, error)
}
