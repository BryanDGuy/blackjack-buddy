# Blackjack Buddy

Blackjack decision assistant implementing basic strategy. Minimal Go HTTP server with an embedded Svelte trainer.

## Features

- Matrix-based strategy lookup with advisor
- Embedded web trainer with round simulation
- Supports hard, soft, doubles, and splits

## Quick Start

**Requirements:** Go 1.25.1+, Make, Node 18+

```bash
git clone https://github.com/bryan/blackjack-buddy.git
cd blackjack-buddy
make build
```

## Usage

```bash
./blackjack-buddy [-port=8080] [-strategy=basic]
```

The binary hosts the HTTP API and embedded trainer. Strategy flag accepts any value from `internal/strategy/strategies` (e.g. `basic`, `coward`).

## Development

- `make ui-build` — rebuild the Svelte trainer only
- `make api-build` — compile the Go server only
- `make run` — build everything then start the server

## Structure

- `api/` — HTTP server, handlers, helpers, embedded assets
- `internal/game` — scenario generation, round resolution
- `internal/strategy` — advisor + strategy matrices
- `web/` — Svelte trainer UI