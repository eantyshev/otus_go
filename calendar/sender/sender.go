package main

import (
	"encoding/json"
	"github.com/eantyshev/otus_go/calendar/pkg/config"
	ent "github.com/eantyshev/otus_go/calendar/pkg/entity"
	"github.com/eantyshev/otus_go/calendar/pkg/logger"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"time"
)

type Sender struct {
	L       *zap.SugaredLogger
	Channel *amqp.Channel
}

func NewSender(l *zap.SugaredLogger, amqpCreds string) (s *Sender) {
	s = &Sender{L: l}
	s.SetupAmqp(amqpCreds)
	return s
}

func (s *Sender) waitConnect(amqpCreds string, timeout, retryPeriod time.Duration) (conn *amqp.Connection, err error) {
	deadline := time.Now().Add(timeout)
	for deadline.After(time.Now()) {
		if conn, err = amqp.Dial(amqpCreds); err == nil {
			return
		}
		s.L.Debug(err)
		time.Sleep(retryPeriod)
	}
	return
}

func (s *Sender) SetupAmqp(amqpCreds string) {
	conn, err := s.waitConnect(
		amqpCreds,
		20*time.Second,
		2*time.Second,
	)
	s.failOnError(err, "failed to connect rabbitmq")
	s.Channel, err = conn.Channel()
	s.failOnError(err, "failed to open a channel")
	_, err = s.Channel.QueueDeclare(
		"calendar.notification", // name
		false,                   // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	s.failOnError(err, "failed to declare a queue")
}
func (s *Sender) failOnError(err error, msg string) {
	if err != nil {
		s.L.Fatalf("%s: %s", msg, err)
	}
}

// PollNotifications fetches all appointments starting on given period
func (s *Sender) ConsumeForever() {
	msgs, err := s.Channel.Consume(
		"calendar.notification",
		"calendar.sender",
		true,
		false,
		false,
		false,
		nil,
	)
	s.failOnError(err, "failed to consume")
	for msg := range msgs {
		ap := &ent.Appointment{}
		err := json.Unmarshal(msg.Body, ap)
		s.failOnError(err, "failed to decode message")
		s.L.Infow("notification received", "owner", ap.Owner, "starts_at", ap.TimeStart)
	}
}

func main() {
	config.Configure()
	amqpDsn := viper.GetString("amqp_dsn")
	s := NewSender(logger.L, amqpDsn)
	s.ConsumeForever()
}
