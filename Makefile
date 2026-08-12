.PHONY: all build clean run ui-build ui-install api-build lint ui-lint api-lint test ui-test api-test help

all: build

ui-install:
	@test -d web/node_modules || npm ci --prefix web

ui-build: ui-install
	npm run build --prefix web

api-build:
	go build -o blackjack-buddy -v ./api

build: ui-build
	$(MAKE) api-build

run: build
	./blackjack-buddy

ui-lint:
	npm run lint --prefix web

api-lint:
	golangci-lint run ./...

lint: ui-lint api-lint

api-test:
	go test ./...

ui-test:
	npm test --prefix web

test: api-test ui-test

clean:
	go clean
	rm -f blackjack-buddy
	rm -rf api/assets

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
