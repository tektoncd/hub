import { when } from 'mobx';
import { getSnapshot } from 'mobx-state-tree';
import { FakeHub } from '../api/testutil';
import { Provider, ProviderStore } from './provider';

const TESTDATA_DIR = `${__dirname}/testdata`;
const api = new FakeHub(TESTDATA_DIR);

describe('Store Object', () => {
  it('can create a provider object', () => {
    const provider = Provider.create({
      name: 'git'
    });

    expect(provider.name).toBe('git');
  });
});

describe('Store functions', () => {
  it('can create a Provider store', (done) => {
    const store = ProviderStore.create({}, { api });
    expect(store.count).toBe(0);
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        expect(store.count).toBe(3);
        expect(store.isLoading).toBe(false);

        expect(getSnapshot(store)).toMatchSnapshot();

        done();
      }
    );
  });
});
