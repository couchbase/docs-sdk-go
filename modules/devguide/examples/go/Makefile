build-student-record:
	cd student-record && go build ./...

build-examples:
	go test ./... -run=^$

build: build-student-record build-examples

lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8
	golangci-lint run -v
