import { UpdateURL } from './updateUrl';

describe('Test UpdateUrl function', () => {
  it('Test UpdateUrl function', () => {
    const val = UpdateURL('', 'rating', 'cli', 'task', 'Tekton', ['cli', 'gke']);
    expect(val).toEqual('sortBy=rating&category=cli&kind=task&catalog=Tekton&tag=cli&tag=gke');
  });
});
