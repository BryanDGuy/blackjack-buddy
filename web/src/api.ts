export interface GameResponse {
  gameId: string;
}

export interface DealResponse {
  playerCards: string[];
  dealerCard: string;
}

export interface MoveRequest {
  move: string;
}

export interface MoveResponse {
  roundState: string;
  activeHand: string[];
  unresolvedHands: string[][];
  resolvedHands: string[][];
  dealerCards: string[];
  outcomes: string[];
}

export interface HintResponse {
  hint: string;
}

export interface ErrorResponse {
  error: string;
  code: string;
}

let gameId: string | null = null;

export async function createGame(): Promise<string> {
  const res = await fetch('/api/game', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' }
  });
  if (!res.ok) {
    const err: ErrorResponse = await res.json();
    throw new Error(err.error || 'Failed to create game');
  }
  const data: GameResponse = await res.json();
  gameId = data.gameId;
  return gameId;
}

export async function loadDeal(): Promise<DealResponse> {
  if (!gameId) {
    await createGame();
  }
  const res = await fetch(`/api/game/${gameId}/deal`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' }
  });
  if (!res.ok) {
    const err: ErrorResponse = await res.json();
    throw new Error(err.error || 'Failed to load deal');
  }
  return res.json();
}

export async function makeMove(move: string): Promise<MoveResponse> {
  if (!gameId) {
    throw new Error('No active game');
  }
  const res = await fetch(`/api/game/${gameId}/move`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ move })
  });
  if (!res.ok) {
    const err: ErrorResponse = await res.json();
    throw new Error(err.error || 'Failed to make move');
  }
  return res.json();
}

export async function getHint(): Promise<string> {
  if (!gameId) {
    return '';
  }
  try {
    const res = await fetch(`/api/game/${gameId}/hint`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' }
    });
    if (!res.ok) {
      return '';
    }
    const data: HintResponse = await res.json();
    return data.hint || '';
  } catch {
    return '';
  }
}
