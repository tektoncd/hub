import { UpdateURL } from './updateUrl';

describe('Test UpdateUrl function', () => {
  it('Test UpdateUrl function', () => {
    const val = UpdateURL('', 'rating', 'cli', 'task', 'Tekton', 'linux/amd64', ['cli', 'gke']);
    expect(val).toEqual(
      'sortBy=rating&category=cli&platform=linux%2Famd64&kind=task&catalog=Tekton&tag=cli&tag=gke'
    );
  });
});
