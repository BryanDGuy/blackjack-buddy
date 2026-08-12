import { afterEach, beforeEach, expect, test, vi } from 'vitest';
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

const deferred = <T>() => {
  let resolve!: (value: T) => void;
  const promise = new Promise<T>(next => { resolve = next; });
  return { promise, resolve };
};

let apps: App[] = [];

const mount = () => {
  const target = document.createElement('div');
  const app = new App({ target });
  apps.push(app);
  return { app, target };
};

beforeEach(() => {
  vi.resetAllMocks();
  document.body.innerHTML = '';
});

afterEach(() => {
  apps.forEach(app => app.$destroy());
  apps = [];
  vi.restoreAllMocks();
  document.body.innerHTML = '';
});

test('disables decisions and keyboard shortcuts until the current hint resolves', async () => {
  const hint = deferred<string>();
  api.loadDeal.mockResolvedValue({ playerCards: ['10', '6'], dealerCard: '10' });
  api.getHint.mockReturnValue(hint.promise);
  const { target } = mount();
  await flush();
  const hit = [...target.querySelectorAll('button')].find(button => button.textContent.startsWith('HIT'))!;
  expect(hit.disabled).toBe(true);
  hit.disabled = false;
  hit.click();
  expect(api.makeMove).not.toHaveBeenCalled();
  window.dispatchEvent(new KeyboardEvent('keydown', { key: 'h' }));
  expect(api.makeMove).not.toHaveBeenCalled();
  hint.resolve('HIT');
  await flush();
  expect(hit.disabled).toBe(false);
});

test('abandons an active round before dealing after a wrong answer', async () => {
  const abandon = deferred<void>();
  api.loadDeal.mockResolvedValue({ playerCards: ['10', '6'], dealerCard: '10' });
  api.getHint.mockResolvedValue('HIT');
  api.makeMove.mockResolvedValue({ roundState: 'active', activeHand: ['10', '6'], unresolvedHands: [], resolvedHands: [], dealerCards: ['10'], outcomes: [] });
  api.abandonRound.mockReturnValue(abandon.promise);
  const { target } = mount();
  await flush();
  const stand = [...target.querySelectorAll('button')].find(button => button.textContent.startsWith('STAND'))!;
  stand.click();
  await flush();
  const next = target.querySelector<HTMLButtonElement>('.next-btn')!;
  next.click();
  await flush();
  expect(api.abandonRound).toHaveBeenCalledOnce();
  expect(api.loadDeal).toHaveBeenCalledTimes(1);
  abandon.resolve();
  await flush();
  expect(api.loadDeal).toHaveBeenCalledTimes(2);
});

test('renders every resolved split hand', async () => {
  api.loadDeal.mockResolvedValue({ playerCards: ['8', '8'], dealerCard: '6' });
  api.getHint.mockResolvedValue('SPLIT');
  api.makeMove.mockResolvedValue({ roundState: 'complete', activeHand: [], unresolvedHands: [], resolvedHands: [['8', '10'], ['8', '9']], dealerCards: ['6', '10', '2'], outcomes: ['Push', 'Lose'] });
  const { target } = mount();
  await flush();
  const split = [...target.querySelectorAll('button')].find(button => button.textContent.startsWith('SPLIT'))!;
  split.click();
  await flush();
  expect(target.querySelectorAll('.resolved-hand')).toHaveLength(2);
});

test('does not reuse a stale hint after an equal-key hand', async () => {
  // A late older request is impossible through the public UI: every hand transition
  // requires its current hint to settle. This verifies the closest case: a settled
  // old hint must not enable an equal-key later hand while its fresh hint is pending.
  const freshHint = deferred<string>();
  api.loadDeal.mockResolvedValue({ playerCards: ['10', '6'], dealerCard: '10' });
  api.getHint.mockResolvedValueOnce('HIT').mockReturnValueOnce(freshHint.promise);
  api.makeMove.mockResolvedValue({ roundState: 'active', activeHand: ['10', '6'], unresolvedHands: [], resolvedHands: [], dealerCards: ['10'], outcomes: [] });
  const { target } = mount();
  await flush();
  const hit = [...target.querySelectorAll('button')].find(button => button.textContent.startsWith('HIT'))!;
  hit.click();
  await flush();
  expect(api.getHint).toHaveBeenCalledTimes(2);
  expect(hit.disabled).toBe(true);
  freshHint.resolve('STAND');
  await flush();
  expect(hit.disabled).toBe(false);
});

test('retries dealing without abandoning again after a post-abandon deal failure', async () => {
  api.loadDeal
    .mockResolvedValueOnce({ playerCards: ['10', '6'], dealerCard: '10' })
    .mockRejectedValueOnce(new Error('deal failed'))
    .mockResolvedValueOnce({ playerCards: ['9', '7'], dealerCard: '10' });
  api.getHint.mockResolvedValue('HIT');
  api.makeMove.mockResolvedValue({ roundState: 'active', activeHand: ['10', '6'], unresolvedHands: [], resolvedHands: [], dealerCards: ['10'], outcomes: [] });
  api.abandonRound.mockResolvedValue(undefined);
  vi.spyOn(window, 'alert').mockImplementation(() => undefined);
  vi.spyOn(console, 'error').mockImplementation(() => undefined);
  const { target } = mount();
  await flush();
  [...target.querySelectorAll('button')].find(button => button.textContent.startsWith('STAND'))!.click();
  await flush();
  const next = target.querySelector<HTMLButtonElement>('.next-btn')!;
  next.click();
  await flush();
  next.click();
  await flush();
  expect(api.abandonRound).toHaveBeenCalledOnce();
  expect(api.loadDeal).toHaveBeenCalledTimes(3);
});

test('locks a round with no hint and exposes Next', async () => {
  api.loadDeal.mockResolvedValue({ playerCards: ['10', '6'], dealerCard: '10' });
  api.getHint.mockResolvedValue('');
  const { target } = mount();
  await flush();
  const hit = [...target.querySelectorAll('button')].find(button => button.textContent.startsWith('HIT'))!;
  const next = target.querySelector<HTMLButtonElement>('.next-btn')!;
  expect(hit.disabled).toBe(true);
  expect(target.textContent).toContain('Hint unavailable. Start next round.');
  expect(next.classList.contains('visible')).toBe(true);
});
