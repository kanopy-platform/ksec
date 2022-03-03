DIST := $(CURDIR)/_dist
VERSION := $(shell grep '^const Version =' pkg/version/version.go | cut -d\" -f2)

.PHONY: all
all: test build

.PHONY: test
test:
	go vet ./...
	go test -cover ./...

.PHONY: build
build:
	go install ./cmd/ksec/

.PHONY: dist
dist:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build -o ksec ./cmd/ksec/
	tar -zcvf $(DIST)/ksec-linux-amd64.tgz ksec README.md LICENSE plugin.yaml
	GOOS=darwin GOARCH=amd64 go build -o ksec ./cmd/ksec/
	tar -zcvf $(DIST)/ksec-macos-amd64.tgz ksec README.md LICENSE plugin.yaml
	GOOS=darwin GOARCH=arm64 go build -o ksec ./cmd/ksec/
	tar -zcvf $(DIST)/ksec-macos-arm64.tgz ksec README.md LICENSE plugin.yaml
	GOOS=windows GOARCH=amd64 go build -o ksec.exe ./cmd/ksec/
	tar -zcvf $(DIST)/ksec-windows-amd64.tgz ksec.exe README.md LICENSE plugin.yaml

.PHONY: clean
clean:
	rm -rf ./_dist ./ksec*
