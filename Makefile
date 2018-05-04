.PHONY: all
all: test build

.PHONY: test
test:
	go vet ./...
	go test -cover ./...

.PHONY: build
build:
	go build -o ksec ./cmd/
