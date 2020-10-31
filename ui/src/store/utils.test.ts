import { assert } from './utils';

describe('Assert function', () => {
  it('Values must be defined', () => {
    const data = new Map();
    data.set(1, 'Tekton-Hub');

    const hub = data.get('1');

    expect(hub).toBe(undefined);
    expect(() => {
      assert(hub);
    }).toThrow(new Error('value must be defined'));
  });
});
