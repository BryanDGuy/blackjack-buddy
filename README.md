# 🃏 Blackjack Buddy

A Blackjack decision assistant that helps you make optimal plays based on your hand and the dealer's up card. Written in Go with a clean command-line interface.

## 🎯 Features

- **Optimal Strategy**: Implements basic strategy for blackjack decision making
- **Interactive CLI**: Easy-to-use command-line interface
- **Simple Card Input**: Easy card input using just ranks (A, K, Q, J, 10, 9, etc.)
- **Hand Analysis**: Evaluates hand values, soft hands, pairs, and special situations
- **Decision Support**: Provides recommendations for Hit, Stand, Double Down, Split, and Surrender

## 🚀 Quick Start

### Prerequisites

- Go 1.25.1 or higher
- Make

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/bryan/blackjack-buddy.git
   cd blackjack-buddy
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Build and run:**
   ```bash
   make run
   ```


## 🎮 Usage

Start the application and follow the prompts:

```bash
make run
```

### Example Session

```
🃏 Welcome to Blackjack Buddy!
A Blackjack decision assistant to help you make optimal plays.

=== New Hand ===
Enter your cards (e.g., 'A 7' or 'A 7'): A 7
Enter dealer's up card (e.g., 'K' or '10'): K

Your hand: [A, 7] (Value: 18)
Dealer's up card: K

🎯 RECOMMENDED ACTION: HIT

💡 Take another card. Your current hand isn't strong enough to stand.
📝 Note: You have a 'soft' hand (ace counted as 11).
```

### Card Input Formats

The application accepts cards using simple rank notation:

- **Face cards**: `A K Q J`
- **Number cards**: `10 9 8 7 6 5 4 3 2`
- **Multiple cards**: `A 7` or `K 5 3`
- **Examples**: `A 7`, `K 5`, `10 6`, `8 8` (for pairs)

## 🔧 Development

**Note**: Make is required for all development tasks. This project uses Make as the single interface for building and running the application.

### Available Commands

All development tasks are handled through the Makefile:

- `make build` - Build the application
- `make clean` - Clean build artifacts
- `make run` - Build and run the application
- `make help` - Show all available targets

### Code Organization

- **Card Package**: Handles card representation (rank-only) and deck creation
- **Hand Package**: Manages hand evaluation, scoring, and special hand types
- **Strategy Package**: Implements basic strategy decision logic
- **CLI**: Interactive command-line interface for user interaction

## 📋 Basic Strategy

This implementation follows standard basic strategy rules:

- **Hard Hands**: Standard hitting/standing rules based on hand value vs dealer up card
- **Soft Hands**: Optimal play for hands containing an ace counted as 11
- **Pairs**: Splitting recommendations for pair hands
- **Double Down**: Strategic doubling on favorable hands
- **Surrender**: Late surrender recommendations for poor hands

## 🚀 Future Enhancements

- [ ] Web interface
- [ ] Advanced counting systems
- [ ] Statistical analysis
- [ ] Training mode with explanations

## 📄 License

This project is open source. Feel free to use and modify as needed.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support

If you have any questions or need help, please open an issue on GitHub.
