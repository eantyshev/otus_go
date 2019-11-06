package appointment

import (
	"github.com/eantyshev/otus_go/calendar/pkg/models"
	"time"
)

type Repository interface {
	Fetch(timeBegin time.Time, num int) ([]*models.Appointment, time.Time, error)
	GetById(id int64) (*models.Appointment, error)
	Update(ap *models.Appointment) error
	Store(ap *models.Appointment) error
	Delete(id int64) error
}
