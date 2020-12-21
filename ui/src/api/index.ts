import axios from 'axios';
import { API_URL } from '../config/constants';
import { ICategory } from '../store/category';
import { IResource, IVersion } from '../store/resource';

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
  resourceVersion(resourceId: number): Promise<IVersion>;
  versionUpdate(versionId: number): Promise<IVersion>;
  authentication(authCode: string): Promise<AuthResponse>;
  readme(rawURL: string): Promise<string>;
  yaml(rawURL: string): Promise<string>;
  getRating(resourceId: number, token: string): Promise<Rating>;
  setRating(resourceId: number, token: string, rating: number): Promise<void | null>;
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
      return axios.get(`${API_URL}/resources`).then((response) => response.data);
    } catch (err) {
      return err.response;
    }
  }

  async authentication(authCode: string) {
    try {
      return axios.post(`${API_URL}/auth/login?code=${authCode}`).then((response) => response.data);
    } catch (err) {
      return err.response;
    }
  }

  async resourceVersion(resourceId: number) {
    try {
      return axios
        .get(`${API_URL}/resource/${resourceId}/versions`)
        .then((response) => response.data);
    } catch (err) {
      return err.response;
    }
  }

  async versionUpdate(versionId: number) {
    try {
      return axios
        .get(`${API_URL}/resource/version/${versionId}`)
        .then((response) => response.data);
    } catch (err) {
      return err.response;
    }
  }

  async readme(rawURL: string) {
    try {
      const URL = rawURL.substring(0, rawURL.lastIndexOf('/') + 1);
      return axios.get(`${URL}/README.md`).then((response) => response.data);
    } catch (err) {
      return err.response;
    }
  }

  async yaml(rawURL: string) {
    try {
      const newLine = '\n';
      return axios.get(`${rawURL}`).then((response) => '```yaml' + newLine + response.data);
    } catch (err) {
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
        .then((response) => response.data);
    } catch (err) {
      return err.response;
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
        .then((response) => response.data);
    } catch (err) {
      return err.response;
    }
  }
}
