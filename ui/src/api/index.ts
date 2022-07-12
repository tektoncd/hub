import axios, { AxiosError, AxiosResponse } from 'axios';
import { API_URL, AUTH_BASE_URL, API_VERSION } from '../config/constants';
import { ICategory } from '../store/category';
import { IResource, IVersion } from '../store/resource';
import { ITokenInfo, IUserProfile } from '../store/auth';
import { ICatalog } from '../store/catalog';
import { IProvider } from '../store/provider';

interface Token {
  token: string;
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
  yaml(
    resourceKey: string,
    version: string
  ): Promise<string | AxiosResponse<unknown, unknown> | undefined>;
  getRating(resourceId: number, token: string): Promise<Rating>;
  setRating(resourceId: number, token: string, rating: number): Promise<void | null>;
  getRefreshToken(refreshToken: string): Promise<ITokenInfo>;
  getAccessToken(accessToken: string): Promise<ITokenInfo>;
  profile(token: string): Promise<IUserProfile>;
  providers(): Promise<IProvider>;
}

export class Hub implements Api {
  async categories() {
    try {
      return axios.get(`${API_URL}/categories`).then((response) => response.data);
    } catch (error) {
      const err = error as AxiosError;
      return err.response;
    }
  }
  async resources() {
    try {
      return axios.get(`${API_URL}/${API_VERSION}/resources`).then((response) => response.data);
    } catch (error) {
      const err = error as AxiosError;
      return err.response;
    }
  }

  async catalogs() {
    try {
      const result = await axios
        .get(`${API_URL}/${API_VERSION}/catalogs`)
        .then((response) => response.data);
      return result;
    } catch (error) {
      const err = error as AxiosError;
      return err.response;
    }
  }

  async authentication(authCode: string) {
    try {
      return axios
        .post(`${AUTH_BASE_URL}/auth/login?code=${authCode}`)
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (error) {
      const err = error as AxiosError;
      return err;
    }
  }

  async resourceVersion(resourceId: number) {
    try {
      return axios
        .get(`${API_URL}/${API_VERSION}/resource/${resourceId}/versions`)
        .then((response) => response.data);
    } catch (error) {
      const err = error as AxiosError;
      return err.response;
    }
  }

  async versionUpdate(versionId: number) {
    try {
      return axios
        .get(`${API_URL}/${API_VERSION}/resource/version/${versionId}`)
        .then((response) => response.data);
    } catch (error) {
      const err = error as AxiosError;
      return err.response;
    }
  }

  async readme(resourceKey: string, version?: string) {
    try {
      const URL = `${API_URL}/${API_VERSION}/resource/${resourceKey}/${version}/readme`;
      return axios.get(URL.toLowerCase()).then((response) => response.data.data.readme);
    } catch (error) {
      const err = error as AxiosError;
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
    } catch (error) {
      const err = error as AxiosError;
      return err.response;
    }
  }

  async getRating(resourceId: number, token: string) {
    try {
      return axios
        .get(`${API_URL}/resource/${resourceId}/rating`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        })
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (error) {
      const err = error as AxiosError;
      return err;
    }
  }

  async setRating(resourceId: number, token: string, rating: number) {
    try {
      return axios
        .put(
          `${API_URL}/resource/${resourceId}/rating`,
          { rating: rating },
          {
            headers: {
              Authorization: `Bearer ${token}`
            }
          }
        )
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (error) {
      const err = error as AxiosError;
      return err;
    }
  }

  async getRefreshToken(refreshToken: string) {
    try {
      return axios({
        method: 'post',
        url: `${AUTH_BASE_URL}/user/refresh/refreshtoken`,
        headers: {
          Authorization: `Bearer ${refreshToken}`
        }
      })
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (error) {
      const err = error as AxiosError;
      return err;
    }
  }

  async getAccessToken(refreshToken: string) {
    try {
      return axios({
        method: 'post',
        url: `${AUTH_BASE_URL}/user/refresh/accesstoken`,
        headers: {
          Authorization: `Bearer ${refreshToken}`
        }
      })
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (error) {
      const err = error as AxiosError;
      return err;
    }
  }

  async profile(token: string) {
    try {
      return axios
        .get(`${AUTH_BASE_URL}/user/info`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        })
        .then((response) => response.data)
        .catch((err) => Promise.reject(err.response));
    } catch (error) {
      const err = error as AxiosError;
      return err;
    }
  }

  async providers() {
    try {
      return axios.get(`${AUTH_BASE_URL}/auth/providers`).then((response) => response.data);
    } catch (error) {
      const err = error as AxiosError;
      return err.response;
    }
  }
}
