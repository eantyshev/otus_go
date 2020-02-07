# Environment
FROM golang:1.13.3

RUN mkdir -p /opt/calendar_tests
WORKDIR /opt/calendar_tests
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go get github.com/cucumber/godog/cmd/godog

CMD godog