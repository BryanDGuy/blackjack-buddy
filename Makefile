# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=blackjack-buddy

.PHONY: all build clean run

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/blackjack-buddy
	
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/blackjack-buddy
	./$(BINARY_NAME)

help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run the application"
	@echo "  help         - Show this help"
