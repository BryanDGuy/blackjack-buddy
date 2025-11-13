# Blackjack Buddy

Blackjack decision assistant implementing basic strategy. Minimal CLI written in Go.

## Features

- Matrix-based basic strategy lookup
- Interactive decision advisor
- Strategy simulation with multiple strategies
- Web-based training mode
- Handles soft hands, pairs, and hard hands

## Quick Start

**Requirements:** Go 1.25.1+, Make

```bash
git clone https://github.com/bryan/blackjack-buddy.git
cd blackjack-buddy
make build
```

## Usage

### Interactive
```bash
./blackjack-buddy
# Strategy (basic/coward): basic
# Your cards: A 7
# Dealer card: K
# [A, 7] | Dealer: K | Action: HIT
```

### Training (Web GUI)
```bash
./blackjack-buddy -train [-strategy=<basic|coward>]
# http://localhost:8080
```

### Simulation
```bash
./blackjack-buddy -sim -strategy=basic -pot=1000 -buyin=10 -rounds=100
# 990.00 | -1.0%
```

**Flags:** `-sim`, `-strategy=<basic|coward>`, `-pot=<amount>`, `-buyin=<amount>`, `-rounds=<count>`, `-verbose`

## Development

```bash
make build   # Build and run: make run
make clean   # Show targets: make help
```

## Structure

`internal/card` | `internal/hand` | `internal/strategy` | `internal/simulator` | `internal/trainer` | `cmd/blackjack-buddy`
