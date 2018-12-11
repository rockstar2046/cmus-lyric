GO_CMD ?=go
GOGET_CMD =${GO_CMD} get -u
BIN_NAME := lyrics


build:
	${GO_CMD} build -o ${BIN_NAME} -v cmd/lyrics.go

install: build
	install -d /usr/local/bin/
	install -m 755 ./${BIN_NAME} /usr/local/bin/${BIN_NAME}

clean:
	rm -f ./${BIN_NAME}*

deps:
	${GOGET_CMD} github.com/gizak/termui

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GO_CMD} build -o ${BIN_NAME} -v cmd/lyrics.go

build-osx:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 ${GO_CMD} build -o ${BIN_NAME}_osx -v cmd/lyrics.go

build-all: build-linux build-osx

help: 
	@sed -nr "s/^([a-z\-]*):(.*)/\1/p" Makefile

.PHONY: build install clean deps build-linux build-osx build-all
