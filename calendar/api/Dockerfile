# Environment
FROM golang:1.12.9 as build-env

RUN mkdir -p /opt/calendar
WORKDIR /opt/calendar
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /opt/service/calendar_api api/main.go

# Release
FROM alpine:latest
COPY --from=build-env /opt/service/calendar_api /bin/calendar_api
ENTRYPOINT ["/bin/calendar_api"]
