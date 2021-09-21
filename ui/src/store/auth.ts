import { flow, getEnv, Instance, types } from 'mobx-state-tree';
import { Api } from '../api';
import { titleCase } from '../common/titlecase';

export const TokenInfo = types.model('TokenInfo', {
  token: types.optional(types.string, ''),
  expiresAt: types.optional(types.number, 0),
  refreshInterval: types.optional(types.string, '')
});

export const UserProfile = types.model('UserProfile', {
  userName: types.optional(types.string, ''),
  name: types.optional(types.string, ''),
  avatarUrl: types.optional(types.string, '')
});

export const Error = types.model('Error', {
  status: types.optional(types.number, 0),
  customMessage: types.optional(types.string, ''),
  serverMessage: types.optional(types.string, '')
});

export type IUserProfile = Instance<typeof UserProfile>;
export type ITokenInfo = Instance<typeof TokenInfo>;
export type IError = Instance<typeof Error>;

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
    profile: types.optional(UserProfile, {}),
    isLoading: true,
    isAuthenticated: false,
    userRating: types.optional(types.number, 0),
    authErr: types.optional(Error, {}),
    ratingErr: types.optional(Error, {}),
    isAuthModalOpen: false
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
    addUserProfile(item: IUserProfile) {
      self.profile.userName = item.userName;
      self.profile.name = item.name;
      self.profile.avatarUrl = item.avatarUrl;
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
    setIsAuthModalOpen(l: boolean) {
      self.isAuthModalOpen = l;
    },
    setErrorMessage(error: IError) {
      switch (error.status) {
        case 400:
          error.customMessage = 'Bad Request';
          break;
        case 401:
          error.customMessage = 'Unauthorized. Please login again!';
          break;
        case 404:
          error.customMessage = 'Requested resource not found';
          break;
        case 500:
          error.customMessage = 'Internal server error. Please try again after some time!';
          break;
        default:
          error.customMessage = '';
      }
      self.authErr = error;
    },
    logout() {
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
    authenticate: flow(function* (authCode: AuthCodeProps, historyBack?: Refresh) {
      try {
        self.setLoading(true);

        const { api } = self;
        const json = yield api.authentication(authCode.code);

        const userDetails = json.data;

        self.addAccessTokenInfo(userDetails.access);
        self.addRefreshTokenInfo(userDetails.refresh);

        self.setIsAuthenticated(true);
        if (historyBack) {
          historyBack();
        }
      } catch (err) {
        if (err === undefined) {
          self.isAuthenticated = false;
        } else {
          const error: IError = {
            status: err.status,
            serverMessage: titleCase(err.data),
            customMessage: ''
          };
          self.setErrorMessage(error);
          self.authErr = error;
        }
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
        if (err === undefined) {
          self.isAuthenticated = false;
        } else {
          const error: IError = {
            status: err.status,
            serverMessage: titleCase(err.data.message),
            customMessage: ''
          };
          self.isAuthenticated = false;
          self.setErrorMessage(error);
          self.ratingErr = error;
        }
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
        const error: IError = {
          status: err.status,
          serverMessage: titleCase(err.data.message),
          customMessage: ''
        };
        self.isAuthenticated = false;
        self.setErrorMessage(error);
        self.ratingErr = error;
      }
      self.setLoading(false);
    }),

    updateRefreshToken: flow(function* () {
      try {
        const { api } = self;

        const refresh = yield api.getRefreshToken(self.refreshTokenInfo.token);
        const newRefreshToken = refresh.data;
        self.addRefreshTokenInfo(newRefreshToken.refresh);
      } catch (err) {
        if (err === undefined) {
          self.isAuthenticated = false;
        } else {
          const error: IError = {
            status: err.status,
            serverMessage: titleCase(err.data),
            customMessage: 'Refresh token has been expired please login again !'
          };
          self.setErrorMessage(error);
          self.setIsAuthenticated(false);
        }
      }
    }),

    updateAccessToken: flow(function* () {
      try {
        const { api } = self;

        const access = yield api.getAccessToken(self.refreshTokenInfo.token);
        const newAccessToken = access.data;
        self.addAccessTokenInfo(newAccessToken.access);
      } catch (err) {
        if (err === undefined) {
          self.isAuthenticated = false;
        } else {
          const error: IError = {
            status: err.status,
            serverMessage: titleCase(err.data),
            customMessage: 'Refresh token has been expired please login again !'
          };
          self.setErrorMessage(error);
          self.setIsAuthenticated(false);
          setTimeout(() => {
            localStorage.clear();
          }, 1000);
        }
      }
    }),

    getProfile: flow(function* () {
      try {
        const { api } = self;

        const access = yield api.profile(self.accessTokenInfo.token);
        const userdata = access.data;
        self.addUserProfile(userdata);
      } catch (err) {
        if (err === undefined) {
          self.isAuthenticated = false;
        } else {
          const error: IError = {
            status: err.status,
            serverMessage: titleCase(err.data),
            customMessage: 'Access token has been expired please login again !'
          };
          self.setErrorMessage(error);
          self.setIsAuthenticated(false);
        }
      }
    })
  }));

export type IAuthStore = Instance<typeof AuthStore>;
