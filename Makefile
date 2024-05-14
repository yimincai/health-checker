GO ?= go
GOFMT ?= gofmt "-s"
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)
GO_FILES := $(shell find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*")
TEST_TAGS ?= ""
GIT_COMMIT_SHA := $(shell git rev-parse HEAD | cut -c 1-8)

default: dev

run:
	@go run main.go

dev: clean
	@CompileDaemon -build="go build -race -o bin/server cmd/main.go" -command="./bin/server" -color=true -graceful-kill=true

clean:
	@go clean
	@-rm -f bin/server

test:
	@gotestsum --junitfile-hide-empty-pkg --format testname

tidy:
	@go mod tidy
	@go fmt ./...

lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u golang.org/x/lint/golint; \
	fi
	@for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

install-tools:
	if [ $(GO_VERSION) -gt 15 ]; then \
		$(GO) install golang.org/x/lint/golint@latest; \
	elif [ $(GO_VERSION) -lt 16 ]; then \
		$(GO) install golang.org/x/lint/golint; \
	fi
