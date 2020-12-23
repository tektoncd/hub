import { flow, getEnv, Instance, types } from 'mobx-state-tree';
import { Api } from '../api';
import { assert } from './utils';

export const TokenInfo = types.model('TokenInfo', {
  token: types.optional(types.string, ''),
  expiresAt: types.optional(types.number, 0),
  refreshInterval: types.optional(types.string, '')
});

export type ITokenInfo = Instance<typeof TokenInfo>;
export interface AuthCodeProps {
  code: string;
}

interface Refresh {
  (): void;
}

export const AuthStore = types
  .model('AuthStore', {
    accessTokenInfo: types.optional(TokenInfo, {}),
    refreshTokenInfo: types.optional(TokenInfo, {}),
    isLoading: true,
    isAuthenticated: false,
    userRating: types.optional(types.number, 0),
    authErr: '',
    ratingErr: ''
  })
  .actions((self) => ({
    addAccessTokenInfo(item: ITokenInfo) {
      self.accessTokenInfo.token = item.token;
      self.accessTokenInfo.expiresAt = item.expiresAt;
      self.accessTokenInfo.refreshInterval = item.refreshInterval;
    },
    addRefreshTokenInfo(item: ITokenInfo) {
      self.refreshTokenInfo.token = item.token;
      self.refreshTokenInfo.expiresAt = item.expiresAt;
      self.refreshTokenInfo.refreshInterval = item.refreshInterval;
    },
    setIsAuthenticated(l: boolean) {
      self.isAuthenticated = l;
    },
    setUserRating(rating: number) {
      self.userRating = rating;
    },
    setLoading(l: boolean) {
      self.isLoading = l;
    },
    onFailure(err: Error) {
      self.authErr = err.toString();
    },
    logout() {
      localStorage.clear();
      self.isAuthenticated = false;
      self.isLoading = false;
      self.userRating = 0;
    }
  }))
  .views((self) => ({
    get api(): Api {
      return getEnv(self).api;
    }
  }))
  .actions((self) => ({
    authenticate: flow(function* (authCode: AuthCodeProps, refresh?: Refresh) {
      try {
        self.setLoading(true);

        const { api } = self;
        const json = yield api.authentication(authCode.code);

        const userDetails = json.data;
        self.addAccessTokenInfo(userDetails.access);
        self.addRefreshTokenInfo(userDetails.refresh);

        self.setIsAuthenticated(true);

        assert(refresh);
        refresh();
      } catch (err) {
        self.authErr = err.toString();
      }
      self.setLoading(false);
    }),

    getRating: flow(function* (resourceId: number) {
      try {
        self.setLoading(true);
        const { api } = self;

        if (self.isAuthenticated) {
          const json = yield api.getRating(resourceId, self.accessTokenInfo.token);
          self.setUserRating(json.rating);
        }
      } catch (err) {
        self.ratingErr = err.toString();
      }
      self.setLoading(false);
    }),

    setRating: flow(function* (resourceId: number, rating: number) {
      try {
        self.setLoading(true);

        const { api } = self;

        if (self.isAuthenticated) {
          yield api.setRating(resourceId, self.accessTokenInfo.token, rating);
          self.setUserRating(rating);
        }
      } catch (err) {
        self.ratingErr = err.toString();
      }
      self.setLoading(false);
    })
  }));
export type IAuthStore = Instance<typeof AuthStore>;
