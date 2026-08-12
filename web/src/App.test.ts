import { beforeEach, expect, test, vi } from 'vitest';
import { tick } from 'svelte';
import App from './App.svelte';

const api = vi.hoisted(() => ({
  loadDeal: vi.fn(),
  makeMove: vi.fn(),
  getHint: vi.fn(),
  abandonRound: vi.fn()
}));
vi.mock('./api', () => api);

const flush = async () => {
  await Promise.resolve();
  await tick();
  await Promise.resolve();
  await tick();
};

beforeEach(() => {
  vi.resetAllMocks();
  document.body.innerHTML = '';
});

test('disables decisions until the current hint resolves', async () => {
  let resolveHint!: (value: string) => void;
  api.loadDeal.mockResolvedValue({ playerCards: ['10', '6'], dealerCard: '10' });
  api.getHint.mockReturnValue(new Promise(resolve => { resolveHint = resolve; }));
  const target = document.createElement('div');
  const app = new App({ target });
  await flush();
  const hit = [...target.querySelectorAll('button')].find(button => button.textContent?.startsWith('HIT'))!;
  expect(hit.disabled).toBe(true);
  resolveHint('HIT');
  await flush();
  expect(hit.disabled).toBe(false);
  app.$destroy();
});

test('abandons an active round before dealing after a wrong answer', async () => {
  api.loadDeal.mockResolvedValue({ playerCards: ['10', '6'], dealerCard: '10' });
  api.getHint.mockResolvedValue('HIT');
  api.makeMove.mockResolvedValue({ roundState: 'active', activeHand: ['10', '6'], unresolvedHands: [], resolvedHands: [], dealerCards: ['10'], outcomes: [] });
  api.abandonRound.mockResolvedValue(undefined);
  const target = document.createElement('div');
  const app = new App({ target });
  await flush();
  const stand = [...target.querySelectorAll('button')].find(button => button.textContent?.startsWith('STAND'))!;
  stand.click();
  await flush();
  const next = target.querySelector<HTMLButtonElement>('.next-btn')!;
  next.click();
  await flush();
  expect(api.abandonRound).toHaveBeenCalledOnce();
  expect(api.loadDeal).toHaveBeenCalledTimes(2);
  app.$destroy();
});

test('renders every resolved split hand', async () => {
  api.loadDeal.mockResolvedValue({ playerCards: ['8', '8'], dealerCard: '6' });
  api.getHint.mockResolvedValue('SPLIT');
  api.makeMove.mockResolvedValue({ roundState: 'complete', activeHand: [], unresolvedHands: [], resolvedHands: [['8', '10'], ['8', '9']], dealerCards: ['6', '10', '2'], outcomes: ['Push', 'Lose'] });
  const target = document.createElement('div');
  const app = new App({ target });
  await flush();
  const split = [...target.querySelectorAll('button')].find(button => button.textContent?.startsWith('SPLIT'))!;
  split.click();
  await flush();
  expect(target.querySelectorAll('.resolved-hand')).toHaveLength(2);
  app.$destroy();
});
