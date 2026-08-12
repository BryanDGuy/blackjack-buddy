# Blackjack Buddy

Blackjack decision assistant implementing basic strategy. Minimal Go HTTP server with an embedded Svelte UI training module.

![Blackjack Buddy UI](web-ui.png)

## Features

- Matrix-based strategy lookup
- Embedded web trainer with round simulation
- Supports hard, soft, doubles, and splits
- Finite multi-deck shoe with card tracking and reshuffling
- Round outcome tracking across normal and split hands
- Six decks, dealer stands on soft 17 (S17), double after split, no surrender, no insurance
- TypeScript-based Svelte UI with component architecture

## Quick Start

**Requirements:** Go 1.25.1+, Make, Node 18+

```bash
git clone https://github.com/bryan/blackjack-buddy.git
cd blackjack-buddy
make build
```

## Usage

```bash
./blackjack-buddy [-port=8080]
```

The binary hosts the HTTP API and embedded trainer UI.

## Development

- `make ui-build` — rebuild the Svelte UI only
- `make api-build` — compile the Go server only
- `make run` — build everything then start the server

## Structure

- `api/` — HTTP server, handlers, helpers, embedded assets
- `internal/` — internal packages intended to be shared amongst layers
- `web/` — TypeScript Svelte trainer UI with components
