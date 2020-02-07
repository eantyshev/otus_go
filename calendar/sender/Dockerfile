# Environment
FROM golang:1.13.3 as build-env

RUN mkdir -p /opt/calendar_sender
WORKDIR /opt/calendar_sender
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /opt/service/calendar_sender sender/sender.go

# Release
FROM alpine:latest
COPY --from=build-env /opt/service/calendar_sender /bin/calendar_sender
RUN mkdir /etc/calendar
COPY config.yaml /etc/calendar/config.yaml
ENTRYPOINT ["/bin/calendar_sender", "--config", "/etc/calendar/config.yaml"]
