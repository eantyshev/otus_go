package main

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	pb "github.com/eantyshev/otus_go/calendar/pkg/adapters/protobuf"
	ent "github.com/eantyshev/otus_go/calendar/pkg/entity"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

var exampleOwner = "Mr.EA"
var exampleAppointment = ent.Appointment{
	Summary:     "summary",
	Description: "description",
	TimeStart:   time.Date(2020, 1, 5, 12, 34, 0, 0, time.UTC),
	TimeEnd:     time.Date(2020, 1, 5, 13, 34, 0, 0, time.UTC),
	Owner:       exampleOwner,
}

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
	foundAp *ent.Appointment
	lastErr error
	newOwner string
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
	p, err := pb.AppointmentToProto(&exampleAppointment)
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
		_, _ = test.cc.DeleteAppointment(test.ctx, &uid)
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
	pbUuid := &pb.UUID{Value: test.createdId}
	pbAp, err := test.cc.GetAppointment(test.ctx, pbUuid)
	if err != nil {
		test.lastErr = err
		return nil
	}
	ap, err := pb.ProtoToAppointment(pbAp.Info, pbAp.Uuid)
	if err != nil {
		return err
	}
	test.foundAp = ap
	return nil
}

func (test *crudTest) iReceiveTheValidProperties() error {
	if test.lastErr != nil {
		return test.lastErr
	}
	// nullify the UUID
	test.foundAp.Uuid = uuid.UUID{}
	if *test.foundAp != exampleAppointment {
		return fmt.Errorf("found: %s, expected: %s", *test.foundAp, exampleAppointment)
	}
	return nil
}

func (test *crudTest) iSendDeleteRequestForGivenId() error {
	pbUuid := &pb.UUID{Value: test.createdId}
	_, err := test.cc.DeleteAppointment(test.ctx, pbUuid)
	return err
}

func (test *crudTest) itFailsWithMessage(arg1 string) error {
	if test.lastErr != nil {
		s := status.Convert(test.lastErr)
		if s.Message() == arg1 {
			return nil
		}
		return fmt.Errorf("error: %s, expected: %s", s.Message(), arg1)
	}
	return fmt.Errorf("expected failure didn't happen")
}

func (test *crudTest) provideAnotherOwner(arg1 string) error {
	test.newOwner = arg1
	return nil
}

func (test *crudTest) iSendUpdateRequestForGivenId() error {
	pbAp, err := pb.AppointmentToProto(&exampleAppointment)
	panicOnErr(err)
	pbAp.Uuid.Value = test.createdId
	pbAp.Info.Owner = test.newOwner

	_, err = test.cc.UpdateAppointment(test.ctx, pbAp)
	return err
}

func (test *crudTest) appointmentsOwnerIs(arg1 string) error {
	if test.foundAp.Owner != arg1 {
		return fmt.Errorf("Owner is unexpected: %s", test.foundAp.Owner)
	}
	return nil
}

func (test *crudTest) otherProperiesAreNot() error {
	// nullify the UUID & owner
	test.foundAp.Uuid = uuid.UUID{}
	test.foundAp.Owner = exampleOwner
	if *test.foundAp != exampleAppointment {
		return fmt.Errorf("found: %s, expected: %s", *test.foundAp, exampleAppointment)
	}
	return nil
}

func CRUDFeatureContext(s *godog.Suite) {
	test := new(crudTest)
	s.BeforeScenario(test.initClient)

	s.Step(`^I send Create request$`, test.iSendCreateRequest)
	s.Step(`^The response is OK$`, theResponseIsOK)
	s.Step(`^the response contains the generated id$`, test.theResponseContainsTheGeneratedId)
	s.Step(`^some appointment is registered$`, test.someAppointmentIsRegistered)
	s.Step(`^I send GetById request for given id$`, test.iSendGetByIdRequestForGivenId)
	s.Step(`^I receive the valid properties$`, test.iReceiveTheValidProperties)
	s.Step(`^I send Delete request for given id$`, test.iSendDeleteRequestForGivenId)
	s.Step(`^it fails with message "([^"]*)"$`, test.itFailsWithMessage)
	s.Step(`^provide another owner "([^"]*)"$`, test.provideAnotherOwner)
	s.Step(`^I send Update request for given id$`, test.iSendUpdateRequestForGivenId)
	s.Step(`^appointment\'s owner is "([^"]*)"$`, test.appointmentsOwnerIs)
	s.Step(`^other properies are not$`, test.otherProperiesAreNot)

	s.AfterScenario(test.tearDownScenario)
	s.AfterFeature(test.stopClient)
}
