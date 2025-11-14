export const KEY_BINDINGS: Record<string, string> = {
  h: 'HIT',
  s: 'STAND',
  d: 'DOUBLE DOWN',
  p: 'SPLIT',
  n: 'NEXT'
};

export const TABLE_HEADERS: string[] = ['', 'A', '10', '9', '8', '7', '6', '5', '4', '3', '2'];

export const HARD_MATRIX: string[][] = [
  ['17+', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'S'],
  ['16', 'H', 'H', 'H', 'H', 'H', 'S', 'S', 'S', 'S', 'S'],
  ['15', 'H', 'H', 'H', 'H', 'H', 'S', 'S', 'S', 'S', 'S'],
  ['14', 'H', 'H', 'H', 'H', 'H', 'S', 'S', 'S', 'S', 'S'],
  ['13', 'H', 'H', 'H', 'H', 'H', 'S', 'S', 'S', 'S', 'S'],
  ['12', 'H', 'H', 'H', 'H', 'H', 'S', 'S', 'S', 'H', 'H'],
  ['11', 'D', 'D', 'D', 'D', 'D', 'D', 'D', 'D', 'D', 'D'],
  ['10', 'H', 'H', 'D', 'D', 'D', 'D', 'D', 'D', 'D', 'D'],
  ['9', 'H', 'H', 'H', 'H', 'H', 'D', 'D', 'D', 'D', 'H'],
  ['8-', 'H', 'H', 'H', 'H', 'H', 'H', 'H', 'H', 'H', 'H']
];

export const SOFT_MATRIX: string[][] = [
  ['A9', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'S'],
  ['A8', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'S'],
  ['A7', 'H', 'H', 'H', 'S', 'S', 'D', 'D', 'D', 'D', 'S'],
  ['A6', 'H', 'H', 'H', 'H', 'H', 'D', 'D', 'D', 'D', 'H'],
  ['A5', 'H', 'H', 'H', 'H', 'H', 'D', 'D', 'D', 'H', 'H'],
  ['A4', 'H', 'H', 'H', 'H', 'H', 'D', 'D', 'D', 'H', 'H'],
  ['A3', 'H', 'H', 'H', 'H', 'H', 'D', 'D', 'H', 'H', 'H'],
  ['A2', 'H', 'H', 'H', 'H', 'H', 'D', 'D', 'H', 'H', 'H']
];

export const PAIR_MATRIX: string[][] = [
  ['AA', 'SP', 'SP', 'SP', 'SP', 'SP', 'SP', 'SP', 'SP', 'SP', 'SP'],
  ['10 10', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'S', 'S'],
  ['9 9', 'S', 'S', 'SP', 'SP', 'S', 'SP', 'SP', 'SP', 'SP', 'SP'],
  ['8 8', 'SP', 'SP', 'SP', 'SP', 'SP', 'SP', 'SP', 'SP', 'SP', 'SP'],
  ['7 7', 'H', 'H', 'H', 'H', 'SP', 'SP', 'SP', 'SP', 'SP', 'SP'],
  ['6 6', 'H', 'H', 'H', 'H', 'H', 'SP', 'SP', 'SP', 'SP', 'SP'],
  ['5 5', 'H', 'H', 'D', 'D', 'D', 'D', 'D', 'D', 'D', 'D'],
  ['4 4', 'H', 'H', 'H', 'H', 'H', 'SP', 'SP', 'H', 'H', 'H'],
  ['3 3', 'H', 'H', 'H', 'H', 'SP', 'SP', 'SP', 'SP', 'SP', 'SP'],
  ['2 2', 'H', 'H', 'H', 'H', 'SP', 'SP', 'SP', 'SP', 'SP', 'SP']
];
