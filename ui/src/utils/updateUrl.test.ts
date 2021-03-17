import { UpdateURL } from './updateUrl';

describe('Test UpdateUrl function', () => {
  it('Test UpdateUrl function', () => {
    const val = UpdateURL('ansible', 'rating', 'cli', 'task', 'Tekton');
    expect(val).toEqual('query=ansible&sortBy=rating&category=cli&kind=task&catalog=Tekton');
  });
});
