export interface SessionResponse {
  sessionId: string;
}

export interface DealResponse {
  playerCards: string[];
  dealerCard: string;
  deckState: { totalCards: number; rankCounts: Record<string, number> };
}

export interface MoveRequest {
  move: string;
}

export interface MoveResponse {
  roundState: string;
  activeHand: string[];
  inactiveHands: string[][];
  completedHands: string[][];
  dealerCards: string[];
  outcomes: string[];
  deckState: { totalCards: number; rankCounts: Record<string, number> };
}

export interface HintResponse {
  hint: string;
}

export interface ErrorResponse {
  error: string;
  code: string;
}

let sessionId: string | null = null;

export async function createSession(): Promise<string> {
  const res = await fetch('/api/session', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' }
  });
  if (!res.ok) {
    const err: ErrorResponse = await res.json();
    throw new Error(err.error || 'Failed to create session');
  }
  const data: SessionResponse = await res.json();
  sessionId = data.sessionId;
  return sessionId;
}

export async function loadDeal(): Promise<DealResponse> {
  if (!sessionId) {
    await createSession();
  }
  const res = await fetch(`/api/session/${sessionId}/deal`, {
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
  if (!sessionId) {
    throw new Error('No active session');
  }
  const res = await fetch(`/api/session/${sessionId}/move`, {
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
  if (!sessionId) {
    return '';
  }
  try {
    const res = await fetch(`/api/session/${sessionId}/hint`, {
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
