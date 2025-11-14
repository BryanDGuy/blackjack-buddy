export const rowDecisions = (row: string[]) => row.slice(1);

export const decisionClass = (val: string) => {
  switch (val) {
    case 'S':
      return 'cell-stand';
    case 'H':
      return 'cell-hit';
    case 'D':
      return 'cell-double';
    case 'SP':
      return 'cell-split';
    default:
      return '';
  }
};

