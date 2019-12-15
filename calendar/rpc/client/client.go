package client

import (
	"context"
	"fmt"
	pb "github.com/eantyshev/otus_go/calendar/pkg/adapters"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"os"
	"time"
)

type CallArgs struct {
	Action string
	RequestJson string
	Uuid string
	Owner string
	Period string
}

func CRUDAppointment(
	conn *grpc.ClientConn,
	timeout time.Duration,
	args CallArgs,
) (err error) {
	var resp proto.Message
	cc := pb.NewCalendarClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ap := &pb.AppointmentInfo{}
	if args.RequestJson != "" {
		if f, err := os.Open(args.RequestJson); err != nil {
			return err
		} else {
			if err := jsonpb.Unmarshal(f, ap); err != nil {
				return err
			}
		}
	}
	switch args.Action {
	case "list":
		req := &pb.ListRequest{Owner: args.Owner}
		switch args.Period {
		case "day":
			req.Period = pb.ListRequest_DAY
		case "week":
			req.Period = pb.ListRequest_WEEK
		case "month":
			req.Period = pb.ListRequest_MONTH
		}
		resp, err = cc.ListAppointments(ctx, req)
	case "create":
		resp, err = cc.CreateAppointment(ctx, ap)
	case "update":
		pbUuid := &pb.UUID{Value: args.Uuid}
		pbAp := &pb.Appointment{
			Uuid: pbUuid,
			Info: ap,
		}
		resp, err = cc.UpdateAppointment(ctx, pbAp)
	case "delete":
		pbUuid := &pb.UUID{Value: args.Uuid}
		resp, err = cc.DeleteAppointment(ctx, pbUuid)
	}
	if err != nil {
		return err
	}
	fmt.Println("response:")
	m := jsonpb.Marshaler{}
	if err := m.Marshal(os.Stdout, resp); err != nil {
		return err
	}

	return nil
}

func RpcCall(
	addrPort string,
	timeout time.Duration,
	args CallArgs,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addrPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()
	return CRUDAppointment(conn, timeout, args)
}
