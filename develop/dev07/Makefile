.PHONY: all test govet golint

all: test govet golint

test: task_test.go
	@go test -v

govet:
	@go vet ./

golint:
	@golint ./