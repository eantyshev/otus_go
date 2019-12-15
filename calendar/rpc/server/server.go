package server

import (
	"context"
	"fmt"
	"github.com/eantyshev/otus_go/calendar/logger"
	pb "github.com/eantyshev/otus_go/calendar/pkg/adapters"
	"github.com/eantyshev/otus_go/calendar/pkg/appointment"
	"github.com/eantyshev/otus_go/calendar/pkg/appointment/repository"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

func proto2Appointment(pbAp *pb.AppointmentInfo, pbUuid *pb.UUID) (ap *appointment.Appointment, err error) {
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
	if pbAp.GetStartsAt() != nil {
		ap.StartsAt, err = ptypes.Timestamp(pbAp.GetStartsAt())
		if err != nil {
			return nil, fmt.Errorf("failed to parse starts_at: %s", err)
		}
	}
	if pbAp.GetDuration() != nil {
		ap.Duration, err = ptypes.Duration(pbAp.GetDuration())
		if err != nil {
			return nil, fmt.Errorf("failed to parse duration: %s", err)
		}
	}
	return ap, nil
}

func appointment2Proto(ap *appointment.Appointment) (pbAp *pb.Appointment, err error) {
	pbAp = &pb.Appointment{
		Uuid: &pb.UUID{Value: ap.Uuid.String()},
		Info: &pb.AppointmentInfo{
			Summary:     ap.Summary,
			Description: ap.Description,
			Owner:       ap.Owner,
			Duration:    ptypes.DurationProto(ap.Duration),
		},
	}
	if pbAp.Info.StartsAt, err = ptypes.TimestampProto(ap.StartsAt); err != nil {
		return nil, err
	}
	return pbAp, nil
}

type calendarService struct {
	repo repository.Repository
	l    *zap.SugaredLogger
}

// implements pb.CalendarService
func (cs *calendarService) CreateAppointment(
	ctx context.Context,
	req *pb.AppointmentInfo,
) (resp *pb.Appointment, err error) {
	var ap *appointment.Appointment
	if ap, err = proto2Appointment(req, nil); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := cs.repo.Store(ap); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if resp, err = appointment2Proto(ap); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (cs *calendarService) UpdateAppointment(
	ctx context.Context,
	req *pb.Appointment,
) (resp *pb.Appointment, err error) {
	var ap *appointment.Appointment
	if ap, err = proto2Appointment(req.Info, req.Uuid); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := cs.repo.Update(ap); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if resp, err = appointment2Proto(ap); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil

}

func (cs *calendarService) DeleteAppointment(
	ctx context.Context,
	pbUuid *pb.UUID,
) (resp *empty.Empty, err error) {
	resp = &empty.Empty{}
	var id uuid.UUID
	if id, err = uuid.Parse(pbUuid.Value); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err = cs.repo.Delete(id); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

func (cs *calendarService) ListAppointments(
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
	if aps, _, err = cs.repo.Fetch(since, -1); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp = &pb.ListResponse{
		Appointments: make([]*pb.Appointment, len(aps)),
	}
	cs.l.Debug("aps len ", len(aps))
	for i, ap := range aps {
		resp.Appointments[i], err = appointment2Proto(ap)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return resp, nil
}

func newCalendarServer() *calendarService {
	return &calendarService{
		repo: repository.NewMapRepo(),
		l:    logger.L,
	}
}

func Server(addrPort string) {
	lis, err := net.Listen("tcp", addrPort)
	if err != nil {
		panic(err)
	}
	logger.L.Debugf("listening at %s", addrPort)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_zap.UnaryServerInterceptor(logger.L.Desugar()),
	))
	pb.RegisterCalendarServer(grpcServer, newCalendarServer())
	panic(grpcServer.Serve(lis))
}
