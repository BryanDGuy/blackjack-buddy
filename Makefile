.PHONY: all build clean run ui-build ui-install api-build lint ui-lint api-lint test ui-test api-test

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
