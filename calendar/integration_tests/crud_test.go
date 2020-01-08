package main

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	pb "github.com/eantyshev/otus_go/calendar/pkg/adapters/protobuf"
	ent "github.com/eantyshev/otus_go/calendar/pkg/entity"
	"google.golang.org/grpc"
	"time"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

type crudTest struct {
	cc        pb.CalendarClient
	conn      *grpc.ClientConn
	ctx       context.Context
	cancel    context.CancelFunc
	createdId string
}

func (test *crudTest) initClient(interface{}) {
	var err error
	// set global deadline to +10min
	deadline := time.Now().Add(10 * time.Minute)
	test.ctx, test.cancel = context.WithDeadline(context.Background(), deadline)
	test.conn, err = grpc.DialContext(test.ctx, "localhost:8888", grpc.WithInsecure(), grpc.WithBlock())
	panicOnErr(err)
	test.cc = pb.NewCalendarClient(test.conn)
}

func (test *crudTest) stopClient(feature *gherkin.Feature) {
	test.cancel()
	panicOnErr(test.conn.Close())
}

func (test *crudTest) iSendCreateRequest() error {

	ap := &ent.Appointment{
		Summary:     "summary",
		Description: "description",
		TimeStart:   time.Date(2020, 1, 5, 12, 34, 0, 0, time.UTC),
		TimeEnd:     time.Date(2020, 1, 5, 13, 34, 0, 0, time.UTC),
		Owner:       "Mr.EA",
	}
	p, err := pb.AppointmentToProto(ap)
	if err != nil {
		return err
	}
	resp, err := test.cc.CreateAppointment(test.ctx, p.Info)
	if err != nil {
		return err
	}
	test.createdId = resp.Value
	return nil
}

func (test *crudTest) tearDownScenario(interface{}, error) {
	if test.createdId != "" {
		uid := pb.UUID{Value: test.createdId}
		_, err := test.cc.DeleteAppointment(test.ctx, &uid)
		panicOnErr(err)
	}
}

func (test *crudTest) theResponseContainsTheGeneratedId() error {
	if test.createdId == "" {
		return fmt.Errorf("created id is not defined")
	}
	return nil
}

func theResponseIsOK() error {
	return nil
}

func (test *crudTest) someAppointmentIsRegistered() error {
	return test.iSendCreateRequest()
}

func (test *crudTest) iSendGetByIdRequestForGivenId() error {
	uid := pb.UUID{Value: test.createdId}
	ap, err := test.cc.GetAppointment(test.ctx, &uid)
	panicOnErr(err)
}

func iReceiveTheValidProperties() error {
	return godog.ErrPending
}

func iSendDeleteRequestForGivenId() error {
	return godog.ErrPending
}

func itReturnsErrNoSuchId() error {
	return godog.ErrPending
}

func provideAnotherOwner(arg1 string) error {
	return godog.ErrPending
}

func iSendUpdateRequestForGivenId() error {
	return godog.ErrPending
}

func appointmentsOwnerIs(arg1 string) error {
	return godog.ErrPending
}

func otherProperiesAreNot() error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	test := new(crudTest)
	s.BeforeScenario(test.initClient)

	s.Step(`^I send Create request$`, test.iSendCreateRequest)
	s.Step(`^The response is OK$`, theResponseIsOK)
	s.Step(`^the response contains the generated id$`, test.theResponseContainsTheGeneratedId)
	s.Step(`^some appointment is registered$`, test.someAppointmentIsRegistered)
	s.Step(`^I send GetById request for given id$`, iSendGetByIdRequestForGivenId)
	s.Step(`^I receive the valid properties$`, iReceiveTheValidProperties)
	s.Step(`^I send Delete request for given id$`, iSendDeleteRequestForGivenId)
	s.Step(`^it returns ErrNoSuchId$`, itReturnsErrNoSuchId)
	s.Step(`^provide another owner "([^"]*)"$`, provideAnotherOwner)
	s.Step(`^I send Update request for given id$`, iSendUpdateRequestForGivenId)
	s.Step(`^appointment\'s owner is "([^"]*)"$`, appointmentsOwnerIs)
	s.Step(`^other properies are not$`, otherProperiesAreNot)

	s.AfterScenario(test.tearDownScenario)
	s.AfterFeature(test.stopClient)
}
