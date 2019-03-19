import { describe, it, expect } from 'test/lib/common';
import { createRefreshIntervalOptions } from '../timepicker/utils';

describe('timepicker utils', () => {
  it('should filter out refresh intervals less than 60s', () => {
    const refreshIntervals = ['5s', '10 seconds', '30s', '60s', '5m', '15m', '1h', '2h', '1d'];
    const expected = [
      {
        text: '60s',
        value: '60s',
      },
      {
        text: '5m',
        value: '5m',
      },
      {
        text: '15m',
        value: '15m',
      },
      {
        text: '1h',
        value: '1h',
      },
      {
        text: '2h',
        value: '2h',
      },
      {
        text: '1d',
        value: '1d',
      },
    ];

    const filtered = createRefreshIntervalOptions(refreshIntervals);
    expect(filtered).toEqual(expected);
  });
});
