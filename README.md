# Blackjack Buddy

Blackjack decision assistant implementing basic strategy. Minimal CLI written in Go.

## Features

- Matrix-based basic strategy lookup
- Interactive decision advisor
- Web-based training mode
- Handles soft hands, pairs, and hard hands

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
# http://localhost:8080
```

The bundled web UI serves from the same address. Strategy flag accepts any value supported by `internal/strategy/strategies` (e.g. `basic`, `coward`).