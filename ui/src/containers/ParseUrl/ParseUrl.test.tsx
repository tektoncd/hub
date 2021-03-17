import { when } from 'mobx';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { root } = createProviderAndStore(api);

describe('ParseUrl component', () => {
  it('it can set url params to resource store', (done) => {
    const { resources } = root;
    when(
      () => {
        return !resources.isLoading;
      },
      () => {
        resources.setURLParams('?/query=ansible');
        expect(resources.urlParams).toBe('?/query=ansible');

        done();
      }
    );
  });
});
