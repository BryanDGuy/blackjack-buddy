# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=blackjack-buddy
WEB_DIR=web
API_DIR=api

.PHONY: all build clean run dev ui-build ui-install api-build help

all: build

ui-install:
	@test -d $(WEB_DIR)/node_modules || npm install --prefix $(WEB_DIR)

ui-build: ui-install
	npm run build --prefix $(WEB_DIR)

api-build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./api

build: ui-build
	$(MAKE) api-build

run: build
	./$(BINARY_NAME)

dev: ui-build
	$(GOCMD) run ./api

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(WEB_DIR)/dist
	rm -rf $(API_DIR)/assets

help:
	@echo "Available targets:"
	@echo "  build  - Build Go server and web bundle"
	@echo "  api-build - Build Go server only"
	@echo "  run    - Build then start the server"
	@echo "  dev    - Rebuild web bundle then run via go run"
	@echo "  clean  - Remove build artifacts"
	@echo "  help   - Show this help"
