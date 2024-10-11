.PHONY: all task test govet golint clean rebuild

all: task test govet golint

task: task.go
	@go build -o myShell

test: task_test.go
	@go test -v

govet:
	@go vet ./...

golint:
	@golint ./...

clean:
	@rm -f myShell

rebuild: clean all