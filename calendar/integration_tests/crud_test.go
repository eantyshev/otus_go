package main

import (
	"github.com/DATA-DOG/godog"
	"github.com/eantyshev/otus_go/calendar/pkg/adapters/protobuf"

	//"github.com/eantyshev/otus_go/calendar/internal/adapters/protobuf"
)

func iSendCreateRequest() error {
	protobuf.Appointment{}
	return godog.ErrPending
}

func theResponseContainsTheGeneratedId() error {
	return godog.ErrPending
}

func appointmentWithIdRegistered(arg1 string) error {
	return godog.ErrPending
}

func iSendGetByIdRequestForId(arg1 string) error {
	return godog.ErrPending
}

func iSendDeleteRequestForId(arg1 string) error {
	return godog.ErrPending
}

func itReturnsErrNoSuchId() error {
	return godog.ErrPending
}

func iSendUpdateRequestForId(arg1 string) error {
	return godog.ErrPending
}

func otherProperiesAreNot() error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I send Create request$`, iSendCreateRequest)
	s.Step(`^the response contains the generated id$`, theResponseContainsTheGeneratedId)
	s.Step(`^appointment with id "([^"]*)" registered$`, appointmentWithIdRegistered)
	s.Step(`^I send GetById request for id "([^"]*)"$`, iSendGetByIdRequestForId)
	s.Step(`^I send Delete request for id "([^"]*)"$`, iSendDeleteRequestForId)
	s.Step(`^it returns ErrNoSuchId$`, itReturnsErrNoSuchId)
	s.Step(`^I send Update request for id "([^"]*)"$`, iSendUpdateRequestForId)
	s.Step(`^other properies are not$`, otherProperiesAreNot)
}

