.PHONY: all test govet golint

all: task test govet golint

task: task.go
	@go run .

test: task_test.go
	@go test -v

govet:
	@go vet ./

golint:
	@golint ./