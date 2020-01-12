package main

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	pb "github.com/eantyshev/otus_go/calendar/pkg/adapters/protobuf"
	ent "github.com/eantyshev/otus_go/calendar/pkg/entity"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"time"
)

type listTest struct {
	cc                          pb.CalendarClient
	conn                        *grpc.ClientConn
	ctx                         context.Context
	timeNow                     time.Time
	createdIds, listedSummaries []string
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

func (test *listTest) iListAppointmentsForForPeriod(owner, period string) (err error) {
	req := &pb.ListRequest{
		Owner: owner,
	}
	switch period {
	case "day": req.Period = pb.ListRequest_DAY
	case "week": req.Period = pb.ListRequest_WEEK
	case "month": req.Period = pb.ListRequest_MONTH
	}
	req.TimeStart, err = ptypes.TimestampProto(test.timeNow)
	panicOnErr(err)
	pbAps, err := test.cc.ListAppointments(test.ctx, req)
	if err != nil {
		return err
	}
	for _, pbAp := range pbAps.Appointments {
		test.listedSummaries = append(test.listedSummaries, pbAp.Info.Summary)
	}
	return nil
}

func (test *listTest) appointmentIsListed(summary string) error {
	var (
		updatedSummaries []string
		isFound bool
	)
	for _, s := range test.listedSummaries {
		if s == summary {
			isFound = true
		} else {
			updatedSummaries = append(updatedSummaries, s)
		}
	}
	test.listedSummaries = updatedSummaries
	if isFound {
		return nil
	}
	return fmt.Errorf("appointment %s is not listed", summary)
}

func (test *listTest) noOtherAppointmentsAreListed() error {
	if len(test.listedSummaries) > 0 {
		return fmt.Errorf("other appointments are listed: %s", test.listedSummaries)
	}
	return nil
}

func ListFeatureContext(s *godog.Suite) {
	test := new(listTest)
	s.BeforeScenario(test.initClient)

	s.Step(`^the wall time is "([^"]*)"$`, test.theWallTimeIs)
	s.Step(`^appointment "([^"]*)" owned by "([^"]*)" starts at "([^"]*)"$`, test.appointmentOwnedByStartsAt)
	s.Step(`^I list appointments for "([^"]*)" for "([^"]*)" period$`, test.iListAppointmentsForForPeriod)
	s.Step(`^appointment "([^"]*)" is listed$`, test.appointmentIsListed)
	s.Step(`^no other appointments are listed$`, test.noOtherAppointmentsAreListed)

	s.AfterScenario(test.tearDownScenario)
	s.AfterFeature(test.stopClient)
}

