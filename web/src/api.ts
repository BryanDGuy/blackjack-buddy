export interface DealResponse {
  playerCards: string[];
  dealerCard: string;
  pot: number;
  bet: number;
  deckState: { totalCards: number; rankCounts: Record<string, number> };
  hint: string;
}

export interface CheckRequest {
  playerCards: string[];
  dealerCard: string;
  decision: string;
  queuedHands: string[][];
  completedHands: string[][];
  pot: number;
  bet: number;
  totalWinnings: number;
}

export interface CheckResponse {
  correct: boolean;
  correctDecision: string;
  outcome?: string;
  roundComplete: boolean;
  restart?: boolean;
  playerCards?: string[];
  dealerCards?: string[];
  queuedHands?: string[][];
  completedHands?: string[][];
  completedOutcomes?: string[];
  pot?: number;
  bet?: number;
  roundWinnings?: number;
  totalWinnings?: number;
  deckState?: { totalCards: number; rankCounts: Record<string, number> };
  hint?: string;
}

export interface HintRequest {
  playerCards: string[];
  dealerCard: string;
}

export interface HintResponse {
  hint: string;
}

export async function loadDeal(skipTrivial: boolean): Promise<DealResponse> {
  const res = await fetch('/api/deal', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ skipTrivial })
  });
  return res.json();
}

export async function checkDecision(req: CheckRequest): Promise<CheckResponse> {
  const res = await fetch('/api/check', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req)
  });
  return res.json();
}

export async function getHint(req: HintRequest): Promise<string> {
  try {
    const res = await fetch('/api/hint', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(req)
    });
    const data: HintResponse = await res.json();
    return data.hint || '';
  } catch {
    return '';
  }
}

