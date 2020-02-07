# Environment
FROM golang:1.13.3 as build-env

RUN mkdir -p /opt/calendar_scheduler
WORKDIR /opt/calendar_scheduler
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /opt/service/calendar_scheduler scheduler/scheduler.go

# Release
FROM alpine:latest
COPY --from=build-env /opt/service/calendar_scheduler /bin/calendar_scheduler
RUN mkdir /etc/calendar
COPY config.yaml /etc/calendar/config.yaml
ENTRYPOINT ["/bin/calendar_scheduler", "--config", "/etc/calendar/config.yaml"]
