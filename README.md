# Blackjack Buddy

Blackjack decision assistant implementing basic strategy. Minimal CLI written in Go.

## Features

- Matrix-based basic strategy lookup
- Supports all decisions: Hit, Stand, Double Down, Split, Surrender
- Handles soft hands, pairs, and hard hands

## Quick Start

**Requirements:** Go 1.25.1+, Make

```bash
git clone https://github.com/bryan/blackjack-buddy.git
cd blackjack-buddy
make run
```

## Usage

```bash
make run
```

**Example:**
```
Your cards: A 7
Dealer card: K
[A, 7] | Dealer: K | Action: HIT
```

**Card input:** `A K Q J 10 9 8 7 6 5 4 3 2`

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
- `internal/strategy` - Basic strategy matrix lookup
- `cmd/blackjack-buddy` - CLI interface

## License

Open source. Use and modify as needed.
