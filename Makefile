GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=midibox
MIDIBOX=pi@midibox

all: test build
build:
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
deploy: build
	ssh $(MIDIBOX) sudo service midibox stop
	scp $(BINARY_NAME) $(MIDOBOX):
	ssh $(MIDIBOX) sudo service midibox start
	rm -f $(BINARY_NAME)
