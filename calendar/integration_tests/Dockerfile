# Environment
FROM golang:1.12.9

RUN mkdir -p /opt/calendar_tests
WORKDIR /opt/calendar_tests
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go get github.com/DATA-DOG/godog/cmd/godog

CMD godog