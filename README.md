# Blackjack Buddy

Blackjack decision assistant implementing basic strategy. Minimal CLI written in Go.

## Features

- Matrix-based basic strategy lookup
- Interactive decision advisor
- Strategy simulation with multiple strategies
- Handles soft hands, pairs, and hard hands

## Quick Start

**Requirements:** Go 1.25.1+, Make

```bash
git clone https://github.com/bryan/blackjack-buddy.git
cd blackjack-buddy
make build
```

## Usage

### Interactive Mode
```bash
./blackjack-buddy
# Your cards: A 7
# Dealer card: K
# [A, 7] | Dealer: K | Action: HIT
```

### Simulation Mode
```bash
./blackjack-buddy -sim -strategy=basic -pot=1000 -buyin=10 -rounds=100
# Output: 990.00 | -1.0%

./blackjack-buddy -sim -strategy=basic -pot=200 -buyin=10 -rounds=5 -verbose
# 1 | 20 vs 17 | +10.00 | 210.00
# 2 | BJ vs 18 | +15.00 | 225.00
```

**Flags:** `-sim`, `-strategy=<basic|coward>`, `-pot=<amount>`, `-buyin=<amount>`, `-rounds=<count>`, `-verbose`

## Development

```bash
make build   # Build
make clean   # Clean artifacts
make run     # Build and run
make help    # Show targets
```

## Structure

- `internal/card` - Card representation
- `internal/hand` - Hand evaluation and scoring
- `internal/strategy` - Strategy matrix lookup
- `internal/simulator` - Game simulation
- `cmd/blackjack-buddy` - CLI interface
