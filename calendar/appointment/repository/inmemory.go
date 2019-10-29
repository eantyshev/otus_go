package repository

import (
	"time"

	"github.com/eantyshev/otus_go/calendar/appointment"
	"github.com/eantyshev/otus_go/calendar/models"
)

type mapRepo map[int64]*models.Appointment

func NewMapRepo() appointment.Repository {
	var r mapRepo = make(map[int64]*models.Appointment)
	return r
}

func (r mapRepo) Fetch(timeFrom time.Time, num int) ([]*models.Appointment, time.Time, error) {
	var cnt = 0
	var aps []*models.Appointment
	var timeEndMax time.Time
	for _, ap := range r {
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
	if _, ok := r[ap.ID]; ok {
		return models.ErrConflictId
	}
	r[ap.ID] = ap
	return nil
}

func (r mapRepo) GetById(id int64) (*models.Appointment, error) {
	if ap, ok := r[id]; ok {
		return ap, nil
	}
	return nil, models.ErrIdNotFound
}

func (r mapRepo) Update(ap *models.Appointment) error {
	r[ap.ID] = ap
	return nil
}

func (r mapRepo) Delete(id int64) error {
	if _, err := r.GetById(id); err != nil {
		return err
	}
	delete(r, id)
	return nil
}
