package repository

import (
	"github.com/eantyshev/otus_go/calendar/pkg/appointment"
	"github.com/google/uuid"
	"sync"
	"time"
)

type mapRepo struct {
	sync.RWMutex
	M map[uuid.UUID]*appointment.Appointment
}

func NewMapRepo() Repository {
	return &mapRepo{M: make(map[uuid.UUID]*appointment.Appointment)}
}

func (r mapRepo) Fetch(timeFrom time.Time, num int) ([]*appointment.Appointment, time.Time, error) {
	var cnt = 0
	var aps = make([]*appointment.Appointment, 0)
	var timeEndMax time.Time
	r.RLock()
	defer r.RUnlock()
	for _, ap := range r.M {
		if ap.StartsAt.After(timeFrom) {
			aps = append(aps, ap)
			cnt++
			timeEnd := ap.StartsAt.Add(ap.Duration)
			if timeEndMax.Before(timeEnd) {
				timeEndMax = timeEnd
			}
		}
		if cnt == num {
			break
		}
	}
	return aps, timeEndMax, nil
}

func (r mapRepo) Store(ap *appointment.Appointment) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.M[ap.Uuid]; ok {
		return appointment.ErrConflictId
	}
	r.M[ap.Uuid] = ap
	return nil
}

func (r mapRepo) GetById(uuid uuid.UUID) (*appointment.Appointment, error) {
	r.RLock()
	defer r.RUnlock()
	if ap, ok := r.M[uuid]; ok {
		return ap, nil
	}
	return nil, appointment.ErrIdNotFound
}

func (r mapRepo) Update(ap *appointment.Appointment) error {
	r.Lock()
	defer r.Unlock()
	if apOrig, ok := r.M[ap.Uuid]; ok {
		*apOrig = *ap
	} else {
		return appointment.ErrIdNotFound
	}
	return nil
}

func (r mapRepo) Delete(uuid uuid.UUID) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.M[uuid]; !ok {
		return appointment.ErrIdNotFound
	}
	delete(r.M, uuid)
	return nil
}
