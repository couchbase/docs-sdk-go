build-student-record:
	cd student-record && go build ./...

build-examples:
	go test ./... -run=^$

build: build-student-record build-examples

set-sdk-version:
	@[ -n "$(SDK_VERSION)" ] || (echo "SDK_VERSION is required (e.g. v2.12.3, master, c84becd)" >&2; exit 1)
	go get github.com/couchbase/gocb/v2@$(SDK_VERSION)
	cd student-record && go get github.com/couchbase/gocb/v2@$(SDK_VERSION)

lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8
	golangci-lint run -v
