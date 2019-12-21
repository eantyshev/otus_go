package usecases

import (
	"context"
	ent "github.com/eantyshev/otus_go/calendar/pkg/entity"
	"github.com/eantyshev/otus_go/calendar/pkg/interfaces"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type UsecasesInterface interface {
	ListOwnerPeriod(ctx context.Context, owner string, timeFrom time.Time, timeTo time.Time) ([]*ent.Appointment, error)
	Create(ctx context.Context, appointment *ent.Appointment) (*uuid.UUID, error)
	Update(ctx context.Context, appointment *ent.Appointment) error
	Delete(ctx context.Context, uid *uuid.UUID) error
	GetById(ctx context.Context, uid *uuid.UUID) (*ent.Appointment, error)
}

// implements UsecasesInterface
type Usecases struct {
	Repo interfaces.Repository
	L    *zap.SugaredLogger
}

func (uc *Usecases) GetById(ctx context.Context, uid *uuid.UUID) (ap *ent.Appointment, err error) {
	return uc.Repo.GetById(ctx, uid)
}

func (uc *Usecases) ListOwnerPeriod(
	ctx context.Context,
	owner string,
	timeFrom time.Time,
	timeTo time.Time,
) (aps []*ent.Appointment, err error) {
	return uc.Repo.ListOwnerPeriod(ctx, owner, timeFrom, timeTo)
}

func (uc *Usecases) Create(ctx context.Context, ap *ent.Appointment) (uid *uuid.UUID, err error) {
	// TODO: implement PG locking at per-owner level
	// First query for other appointments at the same time
	otherAps, err := uc.Repo.ListOwnerPeriod(ctx, ap.Owner, ap.TimeStart, ap.TimeEnd)
	if err != nil {
		return nil, err
	}
	if len(otherAps) > 0 {
		return nil, ent.ErrTimeBusy
	}
	uid, err = uc.Repo.Create(ctx, ap)
	if err != nil {
		return nil, err
	}
	return uid, nil
}

func (uc *Usecases) Update(ctx context.Context, ap *ent.Appointment) error {
	// TODO need to lock concurrent updates
	// Fetch the previous version
	oldAp, err := uc.Repo.GetById(ctx, &ap.Uuid)
	if err != nil {
		return err
	}
	// TODO: use reflection to traverse all parameters
	if ap.Owner == "" {
		ap.Owner = oldAp.Owner
	}
	if ap.Summary == "" {
		ap.Summary = oldAp.Summary
	}
	if ap.Description == "" {
		ap.Description = oldAp.Description
	}
	if ap.TimeStart.IsZero() {
		ap.TimeStart = oldAp.TimeStart
	}
	if ap.TimeEnd.IsZero() {
		ap.TimeEnd = oldAp.TimeEnd
	}
	return uc.Repo.Update(ctx, ap)
}

func (uc *Usecases) Delete(ctx context.Context, uid *uuid.UUID) error {
	return uc.Repo.Delete(ctx, uid)
}
