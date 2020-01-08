
gen:
	PATH=$$PATH:~/go/bin protoc -I api/proto calendar.proto --go_out=plugins=grpc:pkg/adapters/protobuf

api:
	go build -o calendar_api api/main.go
scheduler:
	go build -o calendar_scheduler scheduler/scheduler.go
sender:
	go build -o calendar_sender sender/sender.go
all: api scheduler sender
.PHONY: gen clean api scheduler sender all
