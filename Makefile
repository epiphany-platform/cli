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

APP_REPO := github.com/mkyc/epiphany-wrapper-poc
APP_NAME := e

all: clean get build
run: build run-empty
licenses: clean get build licences-task
test: clean get build test-task

get:
	$(GET) -u -d -v ./...

build:
	$(GET) -d -v ./...
	$(TIDY)
	$(VENDOR)
	$(VET) ./cmd/... ./pkg/...
	$(FMT) ./cmd/... ./pkg/...
	$(BUILD) -x -o $(BUILD_DIR)$(APP_NAME) $(APP_REPO)

clean:
	$(CLEAN) ./cmd/... ./pkg/...
	rm -rf $(BUILD_DIR)

install:
	$(INSTALL) -v ./...

run-empty:
	$(BUILD_DIR)$(APP_NAME)

run-help:
	$(BUILD_DIR)$(APP_NAME) --help

#go get github.com/google/go-licenses first
licences-task:
	rm -rf $(BUILD_DIR)licences
	go-licenses save $(APP_REPO) --save_path="$(BUILD_DIR)licences"
	go-licenses check $(APP_REPO)
	go-licenses csv $(APP_REPO)

test-task:
	$(TEST) -v ./...
