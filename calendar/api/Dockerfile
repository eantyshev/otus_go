# Environment
FROM golang:1.13.3 as build-env

RUN mkdir -p /opt/calendar_api
WORKDIR /opt/calendar_api
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /opt/service/calendar_api api/main.go

# Release
FROM alpine:latest
COPY --from=build-env /opt/service/calendar_api /bin/calendar_api
RUN mkdir /etc/calendar
COPY config.yaml /etc/calendar/config.yaml
ENTRYPOINT ["/bin/calendar_api", "--config", "/etc/calendar/config.yaml"]
