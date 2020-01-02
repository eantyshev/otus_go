package main

import (
	"context"
	"encoding/json"
	"github.com/eantyshev/otus_go/calendar/pkg/config"
	"github.com/eantyshev/otus_go/calendar/pkg/logger"
	"github.com/eantyshev/otus_go/calendar/pkg/adapters/db"
	"github.com/eantyshev/otus_go/calendar/pkg/interfaces"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"time"
)

const (
	Period = 60 * time.Second
)

type Scheduler struct {
	L       *zap.SugaredLogger
	Repo    interfaces.Repository
	Channel *amqp.Channel
}

func NewScheduler(l *zap.SugaredLogger, pgDsn, amqpCreds string) (s *Scheduler) {
	var err error
	s = &Scheduler{L: l}
	s.Repo, err = db.NewPgRepo(pgDsn)
	s.failOnError(err, "Failed to connect PG")
	s.SetupAmqp(amqpCreds)
	return s
}

func (s *Scheduler) failOnError(err error, msg string) {
	if err != nil {
		s.L.Fatalf("%s: %s", msg, err)
	}
}

func (s *Scheduler) SetupAmqp(amqpCreds string) {
	conn, err := amqp.Dial(amqpCreds)
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

// PollNotifications fetches all appointments starting on given period
func (s *Scheduler) PollNotifications(timeFrom, timeTo time.Time) {
	ctx := context.Background()
	aps, err := s.Repo.FetchPeriod(ctx, timeFrom, timeTo)
	s.failOnError(err, "error fetching from db")
	for _, ap := range aps {
		data, _ := json.Marshal(ap)
		err = s.Channel.Publish(
			"",
			"calendar.notification",
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        data,
			})
		s.failOnError(err, "failed to publish event")
	}
}

func (s *Scheduler) PollForever() {
	var timeFrom, timeTo time.Time
	timeTo = time.Now()
	timeFrom = timeTo.Add(-Period)
	pollOnce := func() {
		timeTo = time.Now()
		s.L.Debugf("from: %v, to: %v\n", timeFrom.Format(time.Stamp), timeTo.Format(time.Stamp))
		s.PollNotifications(timeFrom, timeTo)
		timeFrom = timeTo
	}
	timer := time.NewTicker(Period)
	pollOnce()
	for _ = range timer.C {
		pollOnce()
	}
}

func main() {
	config.Configure()
	pgDsn := viper.GetString("pg_dsn")
	amqpDsn := viper.GetString("amqp_dsn")
	s := NewScheduler(logger.L, pgDsn, amqpDsn)
	s.PollForever()
}