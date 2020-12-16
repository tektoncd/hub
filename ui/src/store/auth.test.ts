import { when } from 'mobx';
import { FakeHub } from '../api/testutil';
import { AuthStore, TokenInfo } from './auth';

const TESTDATA_DIR = `${__dirname}/testdata`;
const api = new FakeHub(TESTDATA_DIR);

describe('Store Object', () => {
  it('can create a TokenInfo object', () => {
    const store = TokenInfo.create({
      token: 'abcd',
      expiresAt: 1606280631,
      refreshInterval: '1h0m0s'
    });

    expect(store.refreshInterval).toBe('1h0m0s');
  });
});

describe('Store functions', () => {
  it('can create a auth store', (done) => {
    const store = AuthStore.create({ accessTokenInfo: {}, refreshTokenInfo: {} }, { api });

    expect(store.isLoading).toBe(true);

    const code = {
      code: 'foo'
    };

    store.authenticate(code);
    when(
      () => !store.isLoading,
      () => {
        expect(store.isAuthenticated).toBe(true);

        done();
      }
    );
  });

  it('can test auth store actions', (done) => {
    const store = AuthStore.create({ accessTokenInfo: {}, refreshTokenInfo: {} }, { api });

    const code = {
      code: 'foo'
    };
    store.authenticate(code);

    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        store.setIsAuthenticated(false);
        expect(store.isAuthenticated).toBe(false);

        store.setLoading(false);
        expect(store.isLoading).toBe(false);

        done();
      }
    );
  });
});
