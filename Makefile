MODULE := github.com/kanopy-platform/ksec
CMD_NAME := ksec

VERSION ?= dirty
GIT_COMMIT := $(shell git rev-parse HEAD)
LDFLAGS = "-X '${MODULE}/internal/version.version=${VERSION}' -X '${MODULE}/internal/version.gitCommit=${GIT_COMMIT}'"

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
dist: dist-setup dist-linux dist-darwin dist-windows ## Cross compile binaries into ./dist/

.PHONY: dist-setup
dist-setup:
	mkdir -p ./bin ./dist
	./scripts/plugin.yaml.sh ${VERSION}

.PHONY: dist-build
dist-build: BIN_NAME=${CMD_NAME}-${GOOS}-${GOARCH}-${VERSION}${FILE_EXT}
dist-build:
	go build -ldflags=$(LDFLAGS) -o ./bin/$(BIN_NAME) ./cmd/$(CMD_NAME)/
	tar -zcvf dist/$(BIN_NAME).tgz ./bin/$(BIN_NAME) README.md LICENSE plugin.yaml

.PHONY: dist-amd-arm
dist-amd-arm:
	@$(MAKE) GOARCH=amd64 dist-build
	@$(MAKE) GOARCH=arm64 dist-build

.PHONY: dist-linux
dist-linux:
	@$(MAKE) GOOS=linux dist-amd-arm

.PHONY: dist-darwin
dist-darwin:
	@$(MAKE) GOOS=darwin dist-amd-arm

.PHONY: dist-windows
dist-windows:
	@$(MAKE) GOOS=windows FILE_EXT=.exe dist-amd-arm

.PHONY: clean
clean:
	rm -rf ./bin ./dist plugin.yaml
