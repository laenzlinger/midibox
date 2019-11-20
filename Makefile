GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=midibox

all: test build
build:
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
deploy: build
	scp $(BINARY_NAME) pi@midibox:
	rm -f $(BINARY_NAME)
