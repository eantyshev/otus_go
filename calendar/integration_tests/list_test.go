package main

import (
	"context"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	pb "github.com/eantyshev/otus_go/calendar/pkg/adapters/protobuf"
	ent "github.com/eantyshev/otus_go/calendar/pkg/entity"
	"google.golang.org/grpc"
	"time"
)

type listTest struct {
	cc        pb.CalendarClient
	conn      *grpc.ClientConn
	ctx       context.Context
	timeNow time.Time
	createdIds []string
	lastErr error
	newOwner string
}

func (test *listTest) initClient(interface{}) {
	var err error
	// set global deadline to +10min
	deadline := time.Now().Add(10 * time.Minute)
	test.ctx, _ = context.WithDeadline(context.Background(), deadline)
	test.conn, err = grpc.DialContext(test.ctx, "localhost:8888", grpc.WithInsecure(), grpc.WithBlock())
	panicOnErr(err)
	test.cc = pb.NewCalendarClient(test.conn)
}

func (test *listTest) stopClient(feature *gherkin.Feature) {
	panicOnErr(test.conn.Close())
}

func (test *listTest) tearDownScenario(interface{}, error) {
	for _, createdId := range test.createdIds {
		uid := pb.UUID{Value: createdId}
		_, err := test.cc.DeleteAppointment(test.ctx, &uid)
		panicOnErr(err)
	}
	test.createdIds = []string{}
}

func (test *listTest) theWallTimeIs(arg1 string) (err error) {
	test.timeNow, err = time.Parse("2006-01-02 15:04", arg1)
	return err
}

func (test *listTest) appointmentOwnedByStartsAt(summary, owner, starts_at string) error {
	timeStart, err := time.Parse("2006-01-02 15:04", starts_at)
	if err != nil {
		return err
	}
	ap := ent.Appointment{
		Summary:     summary,
		TimeStart:   timeStart,
		TimeEnd:     timeStart.Add(time.Hour),
		Owner:       owner,
	}
	pbAp, err := pb.AppointmentToProto(&ap)
	if err != nil {
		return err
	}
	resp, err := test.cc.CreateAppointment(test.ctx, pbAp.Info)
	if err != nil {
		return err
	}
	test.createdIds = append(test.createdIds, resp.Value)
	return nil
}

func iListAppointmentsForForPeriod(arg1, arg2 string) error {
	return godog.ErrPending
}

func appointmentIsListed(arg1 string) error {
	return godog.ErrPending
}

func noOtherAppointmentsAreListed() error {
	return godog.ErrPending
}

func ListFeatureContext(s *godog.Suite) {
	test := new(listTest)
	s.BeforeScenario(test.initClient)

	s.Step(`^the wall time is "([^"]*)"$`, test.theWallTimeIs)
	s.Step(`^appointment "([^"]*)" owned by "([^"]*)" starts at "([^"]*)"$`, test.appointmentOwnedByStartsAt)
	s.Step(`^I list appointments for "([^"]*)" for "([^"]*)" period$`, iListAppointmentsForForPeriod)
	s.Step(`^appointment "([^"]*)" is listed$`, appointmentIsListed)
	s.Step(`^no other appointments are listed$`, noOtherAppointmentsAreListed)

	s.AfterScenario(test.tearDownScenario)
	s.AfterFeature(test.stopClient)
}

