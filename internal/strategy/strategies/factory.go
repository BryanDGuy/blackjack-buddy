package strategies

import "github.com/bryan/blackjack-buddy/internal/strategy"

type StrategyType string

const (
	BasicStrategy  StrategyType = "basic"
	CowardStrategy StrategyType = "coward"
)

func CreateStrategy(strategyType StrategyType) strategy.Strategy {
	switch strategyType {
	case BasicStrategy:
		return NewBasic()
	case CowardStrategy:
		return NewCoward()
	default:
		return NewBasic()
	}
}
