package adapters

import (
	"fmt"
	"github.com/eantyshev/otus_go/calendar/internal/entity"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"time"
)

func ProtoToAppointment(pbAp *AppointmentInfo, pbUuid *UUID) (ap *appointment.Appointment, err error) {
	ap = &appointment.Appointment{
		Summary:     pbAp.GetSummary(),
		Description: pbAp.GetDescription(),
		Owner:       pbAp.GetOwner(),
	}
	if pbUuid != nil {
		if ap.Uuid, err = uuid.Parse(pbUuid.Value); err != nil {
			return nil, err
		}
	} else {
		ap.Uuid = uuid.Must(uuid.NewRandom())
	}
	timeOptions := []struct {
		src *timestamp.Timestamp
		dst *time.Time
	}{
		{pbAp.GetTimeStart(), &ap.TimeStart},
		{pbAp.GetTimeEnd(), &ap.TimeEnd},
	}
	for _, x := range timeOptions {
		if x.src != nil {
			ts, err := ptypes.Timestamp(x.src)
			if err != nil {
				return nil, fmt.Errorf("failed to parse time_start: %s", err)
			}
			*x.dst = ts
		}
	}
	return ap, nil
}

func AppointmentToProto(ap *appointment.Appointment) (pbAp *Appointment, err error) {
	pbAp = &Appointment{
		Uuid: &UUID{Value: ap.Uuid.String()},
		Info: &AppointmentInfo{
			Summary:     ap.Summary,
			Description: ap.Description,
			Owner:       ap.Owner,
		},
	}
	if pbAp.Info.TimeStart, err = ptypes.TimestampProto(ap.TimeStart); err != nil {
		return nil, err
	}
	if pbAp.Info.TimeEnd, err = ptypes.TimestampProto(ap.TimeEnd); err != nil {
		return nil, err
	}
	return pbAp, nil
}
