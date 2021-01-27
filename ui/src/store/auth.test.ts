import { when } from 'mobx';
import { FakeHub } from '../api/testutil';
import { AuthStore, TokenInfo, Error, IError } from './auth';

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
  it('can create Error object', () => {
    const error = Error.create({
      status: 400,
      customMessage: 'Test Custom Message',
      serverMessage: 'Test Server Message'
    });

    expect(error.status).toBe(400);
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

        store.logout();
        expect(store.isAuthenticated).toBe(false);
        expect(store.userRating).toBe(0);

        done();
      }
    );
  });

  it('get user rating for buildah resource', (done) => {
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
        store.getRating(13);

        when(
          () => !store.isLoading,
          () => {
            expect(store.userRating).toBe(2);
            done();
          }
        );
      }
    );
  });

  it('can set the rating for buildah resource', (done) => {
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
        store.setRating(13, 3);

        when(
          () => !store.isLoading,
          () => {
            expect(store.userRating).toBe(3);
            done();
          }
        );
      }
    );
  });

  it('can get the custom error message for status code 400', (done) => {
    const store = AuthStore.create({ accessTokenInfo: {}, refreshTokenInfo: {} }, { api });

    const error: IError = {
      status: 400,
      serverMessage: 'Github Login Failed',
      customMessage: ''
    };
    store.setErrorMessage(error);
    expect(error.customMessage).toBe('Bad Request');

    done();
  });

  it('can update refresh and access token', (done) => {
    const store = AuthStore.create({ accessTokenInfo: {}, refreshTokenInfo: {} }, { api });

    const code = {
      code: 'foo'
    };
    store.authenticate(code);
    store.updateAccessToken();
    store.updateRefreshToken();
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        setTimeout(() => {
          expect(store.isAuthenticated).toBe(true);

          expect(store.accessTokenInfo.refreshInterval).toBe('6h0m0s');
          expect(store.refreshTokenInfo.refreshInterval).toBe('24h0m0s');

          done();
        }, 0);
      }
    );
  });
});
