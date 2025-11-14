import { describe, it, expect } from 'vitest';
import { rowDecisions, decisionClass } from './utils';

describe('utils', () => {
  describe('rowDecisions', () => {
    it('returns all elements except the first', () => {
      expect(rowDecisions(['A', 'S', 'H', 'D'])).toEqual(['S', 'H', 'D']);
    });

    it('returns empty array for single element', () => {
      expect(rowDecisions(['A'])).toEqual([]);
    });
  });

  describe('decisionClass', () => {
    it('returns cell-stand for S', () => {
      expect(decisionClass('S')).toBe('cell-stand');
    });

    it('returns cell-hit for H', () => {
      expect(decisionClass('H')).toBe('cell-hit');
    });

    it('returns cell-double for D', () => {
      expect(decisionClass('D')).toBe('cell-double');
    });

    it('returns cell-split for SP', () => {
      expect(decisionClass('SP')).toBe('cell-split');
    });

    it('returns undefined for unknown value', () => {
      expect(decisionClass('X')).toBe('');
    });
  });
});

