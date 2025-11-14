import { loadDeal, checkDecision, getHint, type CheckResponse } from './api';
import { KEY_BINDINGS } from './constants';

export interface GameState {
  skipTrivial: boolean;
  playerCards: string[];
  dealerCard: string;
  dealerCards: string[];
  queuedHands: string[][];
  completedHands: string[][];
  correct: number;
  total: number;
  pot: number;
  bet: number;
  totalWinnings: number;
  roundWinnings: number;
  deckState: { totalCards: number; rankCounts: Record<string, number> };
  hint: string;
  resultText: string;
  resultClass: string;
  outcomeText: string;
  outcomeClass: string;
  nextVisible: boolean;
  busy: boolean;
  locked: boolean;
}

export function formatOutcomeText(result: CheckResponse): string {
  if (result.outcome) {
    let winningsText = '';
    if (result.roundComplete) {
      const sign = result.roundWinnings >= 0 ? '+' : '';
      winningsText = ` (${sign}${result.roundWinnings})`;
    }
    return `Outcome: ${result.outcome}${winningsText}`;
  } else if (result.roundComplete) {
    const sign = result.roundWinnings >= 0 ? '+' : '';
    return `Outcome: (${sign}${result.roundWinnings})`;
  }
  return 'Outcome: -';
}

export async function updateHint(state: GameState): Promise<void> {
  if (!state.playerCards.length || !state.dealerCard || state.locked) {
    return;
  }
  state.hint = await getHint({ playerCards: state.playerCards, dealerCard: state.dealerCard });
}

export async function handleLoadDeal(state: GameState): Promise<void> {
  const data = await loadDeal(state.skipTrivial);
  state.playerCards = data.playerCards;
  state.dealerCard = data.dealerCard;
  state.dealerCards = data.dealerCard ? [data.dealerCard] : [];
  state.queuedHands = [];
  state.completedHands = [];
  if (state.total === 0) state.pot = data.pot;
  state.bet = data.bet;
  state.roundWinnings = 0;
  state.deckState = data.deckState;
  state.hint = data.hint;
  state.resultText = 'Result: -';
  state.resultClass = 'result';
  state.outcomeText = 'Outcome: -';
  state.outcomeClass = 'result outcome-box';
  state.nextVisible = false;
  state.locked = false;
}

export async function decide(state: GameState, decision: string): Promise<void> {
  if (state.locked || state.busy) {
    return;
  }

  state.busy = true;

  try {
    const result = await checkDecision({
      playerCards: state.playerCards,
      dealerCard: state.dealerCard,
      decision,
      queuedHands: state.queuedHands,
      completedHands: state.completedHands,
      pot: state.pot,
      bet: state.bet,
      totalWinnings: state.totalWinnings
    });

    state.total += 1;
    if (result.correct) {
      state.correct += 1;
    }

    state.pot = result.pot;
    state.bet = result.bet;
    state.totalWinnings = result.totalWinnings;
    state.deckState = result.deckState;
    state.hint = result.hint;

    state.resultClass = `result ${result.correct ? 'correct' : 'incorrect'}`;
    state.resultText = `${result.correct ? 'Correct' : 'Incorrect'}: ${result.correctDecision}`;
    state.outcomeClass = 'result outcome-box';
    state.outcomeText = formatOutcomeText(result);

    if (!result.correct || result.restart) {
      state.nextVisible = true;
    }

    state.playerCards = result.playerCards;
    state.dealerCards = result.dealerCards;
    state.queuedHands = result.queuedHands;
    state.completedHands = result.completedHands;

    if (result.roundComplete && state.queuedHands.length === 0 && result.correct && !result.restart) {
      state.nextVisible = true;
    }
    if (!result.roundComplete && result.correct && !result.restart) {
      state.nextVisible = false;
    }
  } finally {
    state.locked = state.nextVisible;
    state.busy = false;
  }
}

export async function startNextRound(state: GameState): Promise<void> {
  if (!state.nextVisible || state.busy) {
    return;
  }
  state.busy = true;
  try {
    await handleLoadDeal(state);
  } finally {
    state.busy = false;
  }
}

export function handleKey(state: GameState, event: KeyboardEvent): void {
  const target = event.target as HTMLElement;
  const tag = target?.tagName;
  if (tag && ['INPUT', 'SELECT', 'TEXTAREA'].includes(tag)) {
    return;
  }
  const action = KEY_BINDINGS[event.key.toLowerCase()];
  if (!action) return;
  event.preventDefault();
  if (action === 'NEXT') {
    void startNextRound(state);
    return;
  }
  if (state.locked || state.busy) {
    return;
  }
  void decide(state, action);
}

