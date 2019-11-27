package repository

import (
	"github.com/eantyshev/otus_go/calendar/pkg/appointment"
	"sync"
	"time"

	"github.com/eantyshev/otus_go/calendar/pkg/models"
)

type mapRepo struct {
	sync.RWMutex
	M map[int64]*models.Appointment
}

func NewMapRepo() appointment.Repository {
	return &mapRepo{M: make(map[int64]*models.Appointment)}
}

func (r mapRepo) Fetch(timeFrom time.Time, num int) ([]*models.Appointment, time.Time, error) {
	var cnt = 0
	var aps = make([]*models.Appointment, 0)
	var timeEndMax time.Time
	for _, ap := range r.M {
		if ap.StartsAt.After(timeFrom) {
			aps = append(aps, ap)
			cnt++
			timeEnd := ap.StartsAt.Add(time.Duration(ap.DurationMinutes) * time.Minute)
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

func (r mapRepo) Store(ap *models.Appointment) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.M[ap.ID]; ok {
		return models.ErrConflictId
	}
	r.M[ap.ID] = ap
	return nil
}

func (r mapRepo) GetById(id int64) (*models.Appointment, error) {
	r.RLock()
	defer r.RUnlock()
	if ap, ok := r.M[id]; ok {
		return ap, nil
	}
	return nil, models.ErrIdNotFound
}

func (r mapRepo) Update(ap *models.Appointment) error {
	r.Lock()
	defer r.Unlock()
	if apOrig, ok := r.M[ap.ID]; ok {
		*apOrig = *ap
	} else {
		return models.ErrIdNotFound
	}
	return nil
}

func (r mapRepo) Delete(id int64) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.M[id]; !ok {
		return models.ErrIdNotFound
	}
	delete(r.M, id)
	return nil
}
