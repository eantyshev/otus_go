## Calendar service ##

Simple educational GRPC service loosely following the quidelines from
from https://github.com/OtusTeam/Go/blob/master/project-calendar.md

*Terminology warning: **event** is called **appointment** in this project*

### Local setup by docker-compose
```shell script
$ docker-compose up --build
```

### GRPC client ###
to explore the API please use Evans tool:
```shell script
evans --repl --host localhost --port 8888 --package adapters --service Calendar api/proto/calendar.proto
```

