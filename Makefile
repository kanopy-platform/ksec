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
	go install

.PHONY: dist
dist:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build -o ksec ./main.go
	tar -zcvf $(DIST)/ksec-linux-$(VERSION).tgz ksec README.md LICENSE plugin.yaml
	GOOS=darwin GOARCH=amd64 go build -o ksec ./main.go
	tar -zcvf $(DIST)/ksec-macos-$(VERSION).tgz ksec README.md LICENSE plugin.yaml
	GOOS=windows GOARCH=amd64 go build -o ksec.exe ./main.go
	tar -zcvf $(DIST)/ksec-windows-$(VERSION).tgz ksec.exe README.md LICENSE plugin.yaml

.PHONY: clean
clean:
	rm -rf ./_dist ./ksec*
