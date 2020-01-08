package server

import (
	"context"
	"encoding/json"
	pb "github.com/eantyshev/otus_go/calendar/pkg/adapters/protobuf"
	"github.com/eantyshev/otus_go/calendar/pkg/entity"
	"github.com/eantyshev/otus_go/calendar/pkg/logger"
	"github.com/eantyshev/otus_go/calendar/pkg/usecases"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type CalendarService struct {
	Usecases usecases.UsecasesInterface
	L        *zap.SugaredLogger
}

//var _ CalendarService = (pb.CalendarServer)(nil)

// implements pb.CalendarService
func (cs *CalendarService) CreateAppointment(
	ctx context.Context,
	req *pb.AppointmentInfo,
) (uuid *pb.UUID, err error) {
	var ap *appointment.Appointment
	if ap, err = pb.ProtoToAppointment(req, nil); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	m := jsonpb.Marshaler{}
	s, _ := m.MarshalToString(req)
	logger.L.Debug("create from proto:", s)
	bs, _ := json.Marshal(ap)
	logger.L.Debug("create from entity:", string(bs))
	uid, err := cs.Usecases.Create(ctx, ap)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UUID{Value: uid.String()}, nil
}

func (cs *CalendarService) UpdateAppointment(
	ctx context.Context,
	req *pb.Appointment,
) (resp *empty.Empty, err error) {
	var ap *appointment.Appointment
	resp = &empty.Empty{}
	if ap, err = pb.ProtoToAppointment(req.Info, req.Uuid); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := cs.Usecases.Update(ctx, ap); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (cs *CalendarService) DeleteAppointment(
	ctx context.Context,
	pbUuid *pb.UUID,
) (resp *empty.Empty, err error) {
	resp = &empty.Empty{}
	var id uuid.UUID
	if id, err = uuid.Parse(pbUuid.Value); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err = cs.Usecases.Delete(ctx, &id); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

func (cs *CalendarService) ListAppointments(
	ctx context.Context, req *pb.ListRequest,
) (resp *pb.ListResponse, err error) {
	var since time.Time
	var aps []*appointment.Appointment
	now := time.Now()
	switch req.Period.String() {
	case "DAY":
		since = now.AddDate(0, 0, -1)
	case "WEEK":
		since = now.AddDate(0, 0, -7)
	case "MONTH":
		since = now.AddDate(0, -1, 0)
	}
	if aps, err = cs.Usecases.ListOwnerPeriod(ctx, req.Owner, since, now); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp = &pb.ListResponse{
		Appointments: make([]*pb.Appointment, len(aps)),
	}
	cs.L.Debug("aps len ", len(aps))
	for i, ap := range aps {
		resp.Appointments[i], err = pb.AppointmentToProto(ap)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return resp, nil
}

func (cs *CalendarService) GetAppointment(
	ctx context.Context, pbUuid *pb.UUID,
) (pbAp *pb.Appointment, err error) {
	uid := uuid.MustParse(pbUuid.Value)
	ap, err := cs.Usecases.GetById(ctx, &uid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if pbAp, err = pb.AppointmentToProto(ap); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return pbAp, nil
}
