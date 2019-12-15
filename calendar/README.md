## Calendar service ##

Simple educational GRPC service loosely following the quidelines from
from https://github.com/OtusTeam/Go/blob/master/project-calendar.md

*Terminology warning: **event** is called **appointment** in this project*

### Simple usage scenario: ###

* Create an appointment (having the wall time equal to 2019-12-15T17:00)
```
$ cat appointment.json
{
  "owner": "eantyshev",
  "summary": "The summary 5",
  "description": "decription 5",
  "starts_at": "2019-11-25T12:34:00Z",
  "duration": "2h30m"
}
$ ./calendar rpc_client create --request-json appointment.json
info	logging configured
response:
{"uuid":{"value":"dbbdcbad-dc94-4f86-a67d-f2b5b344fad8"},"info":{"summary":"The summary 5","description":"decription 5","startsAt":"2019-11-25T12:34:00Z","duration":"9000s","owner":"eantyshev"}}
```

* List the past month' appointments (day and week should yield nothing)
```
$ ./calendar rpc_client list --owner eantyshev --period month
info	logging configured
response:
{"appointments":[{"uuid":{"value":"dbbdcbad-dc94-4f86-a67d-f2b5b344fad8"},"info":{"summary":"The summary 5","description":"decription 5","startsAt":"2019-11-25T12:34:00Z","duration":"9000s","owner":"eantyshev"}}]}
$ ./calendar rpc_client list --owner eantyshev --period day
info	logging configured
response:
{}
```
* Update an existing appointment (shift the date up to the current):
```
$ ./calendar rpc_client update --request-json appointment.json --uuid dbbdcbad-dc94-4f86-a67d-f2b5b344fad8
info	logging configured
response:
{"uuid":{"value":"dbbdcbad-dc94-4f86-a67d-f2b5b344fad8"},"info":{"summary":"The summary 5","description":"decription 5","startsAt":"2019-12-15T12:34:00Z","duration":"9000s","owner":"eantyshev"}}
```
* This time listing for the day shows the appointment