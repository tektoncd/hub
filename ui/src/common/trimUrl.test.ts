import { TrimUrl } from './trimUrl';

describe('Test TitleCase function', () => {
  it('Test titleCase function', () => {
    const nullUrl = TrimUrl('');
    const url = TrimUrl('https://api.hub.tekton.dev/');

    expect(nullUrl).toEqual('');
    expect(url).toEqual('https://api.hub.tekton.dev');
  });
});
