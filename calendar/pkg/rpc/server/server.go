package server

import (
	"context"
	"encoding/json"
	"github.com/eantyshev/otus_go/calendar/logger"
	pb "github.com/eantyshev/otus_go/calendar/pkg/adapters"
	"github.com/eantyshev/otus_go/calendar/pkg/adapters/db"
	"github.com/eantyshev/otus_go/calendar/pkg/entity"
	"github.com/eantyshev/otus_go/calendar/pkg/interfaces"
	"github.com/eantyshev/otus_go/calendar/pkg/usecases"
	"github.com/golang/protobuf/jsonpb"
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

type calendarService struct {
	usecases usecases.UsecasesInterface
	logger   *zap.SugaredLogger
}

// implements pb.CalendarService
func (cs *calendarService) CreateAppointment(
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
	uid, err := cs.usecases.Create(ctx, ap)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UUID{Value: uid.String()}, nil
}

func (cs *calendarService) UpdateAppointment(
	ctx context.Context,
	req *pb.Appointment,
) (resp *empty.Empty, err error) {
	var ap *appointment.Appointment
	resp = &empty.Empty{}
	if ap, err = pb.ProtoToAppointment(req.Info, req.Uuid); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := cs.usecases.Update(ctx, ap); err != nil {
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
	if err = cs.usecases.Delete(ctx, &id); err != nil {
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
	if aps, err = cs.usecases.ListOwnerPeriod(ctx, req.Owner, since, now); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp = &pb.ListResponse{
		Appointments: make([]*pb.Appointment, len(aps)),
	}
	cs.logger.Debug("aps len ", len(aps))
	for i, ap := range aps {
		resp.Appointments[i], err = pb.AppointmentToProto(ap)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return resp, nil
}

func newCalendarServer(repo interfaces.Repository) *calendarService {
	return &calendarService{
		usecases: &usecases.Usecases{Repo: repo, L: logger.L},
		logger:   logger.L,
	}
}

func Server(addrPort string, pgDsn string) {
	repo, err := db.NewPgRepo(pgDsn)
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", addrPort)
	if err != nil {
		panic(err)
	}
	logger.L.Debugf("listening at %s", addrPort)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_zap.UnaryServerInterceptor(logger.L.Desugar()),
	))
	pb.RegisterCalendarServer(grpcServer, newCalendarServer(repo))
	panic(grpcServer.Serve(lis))
}
