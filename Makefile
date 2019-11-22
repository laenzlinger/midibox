GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=midibox
MIDIBOX=pi@midibox

.DEFAULT_GOAL := help

.PHONY: help

all: test build

build: ## buld for Raspberry pi
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(BINARY_NAME) -v

test: ## run unit tests
	$(GOTEST) -v ./...

clean: ## clean all temporary files
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

deploy: build ## deploy to the midibox
	ssh $(MIDIBOX) sudo service midibox stop
	scp $(BINARY_NAME) $(MIDOBOX):
	ssh $(MIDIBOX) sudo service midibox start
	rm -f $(BINARY_NAME)


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'