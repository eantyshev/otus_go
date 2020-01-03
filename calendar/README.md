## Calendar service ##

Simple educational GRPC service loosely following the quidelines from
from https://github.com/OtusTeam/Go/blob/master/project-calendar.md

*Terminology warning: **event** is called **appointment** in this project*

### GRPC client ###
to explore the API please use Evans tool:
```shell script
evans --repl --host localhost --port 50051 --package adapters --service Calendar api/proto/calendar.proto
```

## Setup
### Build binaries
```shell script
$ make all
go build -o calendar_api api/main.go
go build -o calendar_scheduler scheduler/scheduler.go
go build -o calendar_sender sender/sender.go
```
### Start processes
```shell script
$ ./calendar_api -c config.yaml
$ ./calendar_scheduler -c config.yaml
$ ./calendar_sender -c config.yaml
```

### Create an event and watch for notification
```shell script
debug	create from entity:{"uuid":"9549342b-a21e-47b6-b53e-d9d99fb4e8c2","summary":"summary","description":"descriprion","time_start":"2020-01-03T16:59:46Z","time_end":"2020-01-03T17:01:26Z","owner":"ea"}
info	finished unary call with code OK	{"grpc.start_time": "2020-01-03T19:59:03+03:00", "system": "grpc", "span.kind": "server", "grpc.service": "adapters.Calendar", "grpc.method": "CreateAppointment", "grpc.code": "OK", "grpc.time_ms": 12.767999649047852}
debug	from: 2020-01-03T16:58:07Z, to: 2020-01-03T16:59:07Z

debug	from: 2020-01-03T16:59:07Z, to: 2020-01-03T17:00:07Z

info	publish notification	{"time_start": "2020-01-03T16:59:46.000Z", "owner": "ea"}
info	notification received	{"owner": "ea", "starts_at": "2020-01-03T16:59:46.000Z"}

```
