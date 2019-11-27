package adapters


import (
	"fmt"
	"github.com/eantyshev/otus_go/calendar/pkg/models"
	"net/http"
	"strconv"
	"time"
)


func NewFormAppointment(r *http.Request) (ap *models.Appointment, err error) {
	var (
		id int64
		startsAt time.Time
		durationMinutes uint64
		isRegular bool
	)
	if err = r.ParseForm(); err != nil {
		return nil, err
	}
	if formId := r.FormValue("id"); formId == "" {
		return nil, fmt.Errorf("id is missing or empty")
	} else {
		if id, err = strconv.ParseInt(formId, 10, 64); err != nil {
			return nil, err
		}
	}
	if r.FormValue("starts_at") != "" {
		if startsAt, err = time.Parse(time.RFC3339, r.FormValue("starts_at")); err != nil {
			return nil, err
		}
	}
	if formMinutes := r.FormValue("duration_minutes"); formMinutes != "" {
		if durationMinutes, err = strconv.ParseUint(formMinutes, 10, 16); err != nil {
			return nil, err
		}
	}
	if formRegular := r.FormValue("is_regular"); formRegular != "" {
		if isRegular, err = strconv.ParseBool(formRegular); err != nil {
			return nil, err
		}
	}
	ap = &models.Appointment{
		ID:              id,
		Summary:         r.FormValue("summary"),
		Description:     r.FormValue("description"),
		StartsAt:        startsAt,
		DurationMinutes: uint16(durationMinutes),
		IsRegular:       isRegular,
		DaysOfWeek:      nil,
	}
	return ap, nil
}
