package appointment

import (
	"github.com/eantyshev/otus_go/calendar/models"
)

type Repository interface {
	//    Fetch(cursor string, num int) ([]*models.Appointment, string, error)
	GetById(id int64) (*models.Appointment, error)
	//    Update(ap *models.Appointment) error
	Store(ap *models.Appointment) error
	//    Delete(id int64) error
}
