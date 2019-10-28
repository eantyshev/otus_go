package repository

import (
    "github.com/eantyshev/otus_go/calendar/models"
    "github.com/eantyshev/otus_go/calendar/appointment"
)


type mapRepo map[int64]*models.Appointment

func NewMapRepo() appointment.Repository {
    var r mapRepo = make(map[int64]*models.Appointment)
	return r
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


