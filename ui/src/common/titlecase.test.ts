import { titleCase } from './titlecase';

describe('Test TitleCase function', () => {
  it('Test titleCase function', () => {
    const val = titleCase('test value');
    expect(val).toEqual('Test Value');
  });
});
