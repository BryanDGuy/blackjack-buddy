# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=blackjack-buddy
WEB_DIR=web
API_DIR=api

.PHONY: all build clean run ui-build ui-install api-build lint ui-lint api-lint test ui-test api-test help

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

ui-lint:
	npm run lint --prefix $(WEB_DIR)

api-lint:
	golangci-lint run ./...

lint: ui-lint
	$(MAKE) api-lint

api-test:
	$(GOCMD) test ./...

ui-test:
	npm test --prefix $(WEB_DIR)

test: api-test
	$(MAKE) ui-test

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(WEB_DIR)/dist
	rm -rf $(API_DIR)/assets

help:
	@echo "Available targets:"
	@echo "  build  - Build Go server and web bundle"
	@echo "  api-build - Build Go server only"
	@echo "  ui-build - Build Svelte UI only"
	@echo "  run    - Build then start the server"
	@echo "  lint   - Run static analysis (Go and UI)"
	@echo "  api-lint - Run Go static analysis"
	@echo "  ui-lint - Run UI static analysis"
	@echo "  test   - Run tests (Go and UI)"
	@echo "  api-test - Run Go tests"
	@echo "  ui-test - Run UI tests"
	@echo "  clean  - Remove build artifacts"
	@echo "  help   - Show this help"
