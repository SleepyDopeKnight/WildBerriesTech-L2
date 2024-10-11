.PHONY: all task test govet golint clean rebuild

all: task test govet golint

task: task.go
	@go build -o mySort

test: task_test.go
	@go test -v

govet:
	@go vet ./...

golint:
	@golint ./...

clean:
	@rm -f mySort

rebuild: clean all