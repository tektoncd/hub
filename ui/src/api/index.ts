import axios from 'axios';
import { API_URL, API_VERSION } from '../config/constants';
import { ICategory } from '../store/category';
import { IResource, IVersion } from '../store/resource';
import { ITokenInfo, IUserProfile } from '../store/auth';
import { ICatalog } from '../store/catalog';
import { IProvider } from '../store/provider';

interface Token {
  refreshInterval: string;
  expiresAt: number;
}

interface TokenData {
  access: Token;
  refresh: Token;
}

export interface AuthResponse {
  data: TokenData;
}

export interface AccessTokenResponse {
  data: string;
}

export interface LogoutResponse {
  data: boolean;
}

export interface Rating {
  rating: number;
}

export interface Api {
  categories(): Promise<ICategory>;
  resources(): Promise<IResource>;
  catalogs(): Promise<ICatalog>;
  resourceVersion(resourceId: number): Promise<IVersion>;
  versionUpdate(versionId: number): Promise<IVersion>;
  authentication(authCode: string): Promise<AuthResponse>;
  readme(resourceKey: string, version: string): Promise<string>;
  yaml(resourceKey: string, version: string): Promise<string>;
  getRating(resourceId: number): Promise<Rating>;
  setRating(resourceId: number, rating: number): Promise<void | null>;
  getRefreshToken(): Promise<ITokenInfo>;
  getAccessToken(): Promise<ITokenInfo>;
  profile(): Promise<IUserProfile>;
  providers(): Promise<IProvider>;
  logout(): Promise<LogoutResponse>;
  accessToken(): Promise<AccessTokenResponse>;
}

export class Hub implements Api {
  async categories() {
    try {
      return axios.get(`${API_URL}/categories`).then((response) => response.data);
    } catch (err) {
      return err.response;
    }
  }

  async resources() {
    try {
      return axios.get(`${API_URL}/${API_VERSION}/resources`).then((response) => response.data);
    } catch (err) {
      return err.response;
    }
  }

  async catalogs() {
    try {
      return axios.get(`${API_URL}/${API_VERSION}/catalogs`).then((response) => response.data);
    } catch (err) {
      return err.response;
    }
  }

  async authentication(authCode: string) {
    try {
      return axios
        .post(`/auth/login?code=${authCode}`)
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (err) {
      return err;
    }
  }

  async resourceVersion(resourceId: number) {
    try {
      return axios
        .get(`${API_URL}/${API_VERSION}/resource/${resourceId}/versions`)
        .then((response) => response.data);
    } catch (err) {
      return err.response;
    }
  }

  async versionUpdate(versionId: number) {
    try {
      return axios
        .get(`${API_URL}/${API_VERSION}/resource/version/${versionId}`)
        .then((response) => response.data);
    } catch (err) {
      return err.response;
    }
  }

  async readme(resourceKey: string, version?: string) {
    try {
      const URL = `${API_URL}/${API_VERSION}/resource/${resourceKey}/${version}/readme`;
      return axios.get(URL.toLowerCase()).then((response) => response.data.data.readme);
    } catch (err) {
      return err.response;
    }
  }

  async yaml(resourceKey: string, version?: string) {
    try {
      const newLine = '\n';
      const URL = `${API_URL}/${API_VERSION}/resource/${resourceKey}/${version}/yaml`;
      return axios
        .get(URL.toLowerCase())
        .then((response) => '```yaml' + newLine + response.data.data.yaml);
    } catch (err) {
      return err.response;
    }
  }

  async getRating(resourceId: number) {
    try {
      return axios
        .get(`/resource/${resourceId}/rating`)
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (err) {
      return err;
    }
  }

  async setRating(resourceId: number, rating: number) {
    try {
      return axios
        .put(`/resource/${resourceId}/rating`, { rating: rating })
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (err) {
      return err;
    }
  }

  async getRefreshToken() {
    try {
      return axios({
        method: 'post',
        url: `/user/refresh/refreshtoken`
      })
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (err) {
      return err;
    }
  }

  async getAccessToken() {
    try {
      return axios({
        method: 'post',
        url: `/user/refresh/accesstoken`
      })
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (err) {
      return err;
    }
  }

  async profile() {
    try {
      return axios
        .get(`/user/info`)
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (err) {
      return err;
    }
  }

  async providers() {
    try {
      return axios.get(`/auth/providers`).then((response) => response.data);
    } catch (err) {
      return err.response;
    }
  }

  async logout() {
    try {
      return axios
        .post(`/user/logout `)
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (err) {
      return err;
    }
  }

  async accessToken() {
    try {
      return axios
        .post(`/user/accesstoken`)
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (err) {
      return err;
    }
  }
}
