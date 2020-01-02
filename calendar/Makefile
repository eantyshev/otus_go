
gen:
	PATH=$$PATH:~/go/bin protoc -I rpc/proto calendar.proto --go_out=plugins=grpc:pkg/adapters
build:
	go build -o calendar
clean:
	go clean -v
.PHONY: gen clean build
