CMD := go
GET := $(CMD) get
BUILD := $(CMD) build
VET := $(CMD) vet
FMT := $(CMD) fmt
CLEAN := $(CMD) clean
INSTALL := $(CMD) install
VENDOR := $(CMD) mod vendor
TIDY := $(CMD) mod tidy
TEST := $(CMD) test

BUILD_DIR := ./output/

APP_REPO := github.com/epiphany-platform/cli
APP_NAME := e

.PHONY: all licenses test clean build get install test-task clean-task get-task get-update-task janitor-task build-task licences-task install-task

all: janitor-task clean-task get-task janitor-task build-task test-task
licenses: licences-task
test: get-task janitor-task test-task
clean: clean-task
build: get-task janitor-task build-task
get: get-update-task
install: install-task
pipeline-test: test-task

test-task:
	$(TEST) -race -v ./...

clean-task:
	$(CLEAN) -x ./cmd/... ./pkg/... ./internal/...
	rm -rf $(BUILD_DIR)

get-task:
	$(GET) -d -v ./...

get-update-task:
	$(GET) -u -d -v ./...

janitor-task:
	$(TIDY)
	$(VENDOR)
	$(VET) ./cmd/... ./pkg/... ./internal/...
	$(FMT) ./cmd/... ./pkg/... ./internal/...
	goimports -l -w ./cmd/ ./pkg/ ./internal/

build-task:
	$(BUILD) -x -o $(BUILD_DIR)$(APP_NAME) $(APP_REPO)

#go get github.com/google/go-licenses first
licences-task:
	rm -rf $(BUILD_DIR)licences
	go-licenses save $(APP_REPO) --save_path="$(BUILD_DIR)licences"
	go-licenses check $(APP_REPO)
	go-licenses csv $(APP_REPO)

install-task:
	$(INSTALL) -v ./...
