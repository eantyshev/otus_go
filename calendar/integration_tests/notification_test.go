package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cucumber/godog/gherkin"
	pb "github.com/eantyshev/otus_go/calendar/pkg/adapters/protobuf"
	ent "github.com/eantyshev/otus_go/calendar/pkg/entity"
	"google.golang.org/grpc"

	"github.com/cucumber/godog"
	"github.com/streadway/amqp"
)

var amqpDSN = os.Getenv("CALENDAR_AMQP_DSN")

func init() {
	if amqpDSN == "" {
		amqpDSN = "amqp://guest:guest@localhost:5672/"
	}
}

const (
	amqpQueue    = "calendar.notification"
	amqpConsumer = "calendar.sender"
)

type grpcTest struct {
	cc        pb.CalendarClient
	conn      *grpc.ClientConn
	ctx       context.Context
	timeNow   time.Time
	createdId string
}

func (test *grpcTest) initClient(interface{}) {
	var err error
	// set global deadline to +10min
	deadline := time.Now().Add(10 * time.Minute)
	test.ctx, _ = context.WithDeadline(context.Background(), deadline)
	test.conn, err = grpc.DialContext(test.ctx, "grpc_api:50051", grpc.WithInsecure(), grpc.WithBlock())
	panicOnErr(err)
	test.cc = pb.NewCalendarClient(test.conn)
}

func (test *grpcTest) appointmentHasStartTimeAtNow() error {
	p, err := pb.AppointmentToProto(&ent.Appointment{
		Summary:     "summary",
		Description: "descr.",
		TimeStart:   time.Now(),
		TimeEnd:     time.Now().Add(time.Hour),
		Owner:       "owner",
	})
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

func (test *grpcTest) stopClient(feature *gherkin.Feature) {
	if test.conn != nil {
		panicOnErr(test.conn.Close())
	}
}

func (test *grpcTest) tearDownScenario(interface{}, error) {
	if test.createdId != "" {
		uid := pb.UUID{Value: test.createdId}
		_, err := test.cc.DeleteAppointment(test.ctx, &uid)
		panicOnErr(err)
	}
}

type notifyTest struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	messages      []string
	messagesMutex *sync.RWMutex
	stopSignal    chan struct{}
}

func (test *notifyTest) startConsuming(interface{}) {

	test.messagesMutex = new(sync.RWMutex)
	test.stopSignal = make(chan struct{})

	var err error

	test.conn, err = amqp.DialConfig(
		amqpDSN,
		amqp.Config{
			Heartbeat: 10 * time.Second,
			Locale:    "en_US",
			Dial:      amqp.DefaultDial(60 * time.Second),
		})
	panicOnErr(err)

	test.ch, err = test.conn.Channel()
	panicOnErr(err)

	events, err := test.ch.Consume(amqpQueue, amqpConsumer,
		true, false, false, false, nil)
	panicOnErr(err)

	go func(stop <-chan struct{}) {
		for {
			select {
			case <-stop:
				return
			case event := <-events:
				ap := &ent.Appointment{}
				err := json.Unmarshal(event.Body, ap)
				panicOnErr(err)
				test.messagesMutex.Lock()
				test.messages = append(test.messages, ap.Summary)
				test.messagesMutex.Unlock()
			}
		}
	}(test.stopSignal)
}

func (test *notifyTest) stopConsuming(interface{}, error) {
	test.stopSignal <- struct{}{}

	panicOnErr(test.ch.Close())
	panicOnErr(test.conn.Close())
	test.messages = nil
}

func (test *notifyTest) notificationIsReceivedWithinSeconds(timeoutSec int) error {
	time.Sleep(time.Duration(timeoutSec) * time.Second) // На всякий случай ждём обработки евента

	test.messagesMutex.RLock()
	defer test.messagesMutex.RUnlock()

	for _, msg := range test.messages {
		if msg == "summary" {
			return nil
		}
	}
	return fmt.Errorf("notification wasn't received: %s", test.messages)
}

func NotifyFeatureContext(s *godog.Suite) {
	ntest := new(notifyTest)
	gtest := new(grpcTest)
	s.BeforeScenario(gtest.initClient)
	s.BeforeScenario(ntest.startConsuming)

	s.Step(`^appointment has start time at now$`, gtest.appointmentHasStartTimeAtNow)
	s.Step(`^notification is received within (\d+) seconds$`, ntest.notificationIsReceivedWithinSeconds)

	s.AfterScenario(ntest.stopConsuming)
	s.AfterScenario(gtest.tearDownScenario)
	s.AfterFeature(gtest.stopClient)
}
