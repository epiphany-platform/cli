CMD := go
GET := $(CMD) get
BUILD := $(CMD) build
VET := $(CMD) vet
FMT := $(CMD) fmt
CLEAN := $(CMD) clean
INSTALL := $(CMD) install
VENDOR := $(CMD) mod vendor
TIDY := $(CMD) mod tidy

BUILD_DIR := ./output/

APP_REPO := github.com/mkyc/epiphany-wrapper-poc
APP_NAME := e

all: clean get build
run: build run-empty

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