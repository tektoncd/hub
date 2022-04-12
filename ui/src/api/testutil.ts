import * as fs from 'fs';
import { ICategory } from '../store/category';
import { IResource, IVersion } from '../store/resource';
import { ITokenInfo, IUserProfile } from '../store/auth';
import { Api, AuthResponse } from './';
import { ICatalog } from '../store/catalog';
import { IProvider } from '../store/provider';

export class FakeHub implements Api {
  dataDir: string;

  constructor(dataDir: string) {
    this.dataDir = dataDir;
  }

  async categories() {
    const data = `${this.dataDir}/categories.json`;
    const ret = () => JSON.parse(fs.readFileSync(data).toString());
    return new Promise<ICategory>((resolve) => {
      setTimeout(() => resolve(ret()), 1000);
    });
  }

  async resources() {
    const data = `${this.dataDir}/resources.json`;

    const ret = () => JSON.parse(fs.readFileSync(data).toString());
    return new Promise<IResource>((resolve) => {
      setTimeout(() => resolve(ret()), 1000);
    });
  }

  async catalogs() {
    const data = `${this.dataDir}/catalogs.json`;

    const ret = () => JSON.parse(fs.readFileSync(data).toString());
    return new Promise<ICatalog>((resolve) => {
      setTimeout(() => resolve(ret()), 1000);
    });
  }

  async resourceVersion(resourceId: number) {
    const data = `${this.dataDir}/resource_versions_${resourceId}.json`;

    const ret = () => JSON.parse(fs.readFileSync(data).toString());
    return new Promise<IVersion>((resolve) => {
      setTimeout(() => resolve(ret()), 1000);
    });
  }

  async versionUpdate(versionId: number) {
    const data = `${this.dataDir}/version_details_${versionId}.json`;

    const ret = () => JSON.parse(fs.readFileSync(data).toString());
    return new Promise<IVersion>((resolve) => {
      setTimeout(() => resolve(ret()), 1000);
    });
  }

  async authentication() {
    const data = `${this.dataDir}/authentication.json`;

    const ret = () => JSON.parse(fs.readFileSync(data).toString());
    return new Promise<AuthResponse>((resolve) => {
      setTimeout(() => resolve(ret()), 1000);
    });
  }

  async readme(resourceKey: string) {
    const splitRawUrl = resourceKey.split('/');
    const resourceName = splitRawUrl[splitRawUrl.length - 1];
    const data = `${this.dataDir}/${resourceName}-Readme.md`;

    const ret = () => fs.readFileSync(data).toString();
    return new Promise<string>((resolve) => {
      setTimeout(() => resolve(ret()), 1000);
    });
  }

  async yaml(resourceKey: string) {
    const splitRawUrl = resourceKey.split('/');
    const resourceName = splitRawUrl[splitRawUrl.length - 1];
    const data = `${this.dataDir}/${resourceName}.yaml`;

    const ret = () => fs.readFileSync(data).toString();
    return new Promise<string>((resolve) => {
      setTimeout(() => resolve(ret()), 1000);
    });
  }

  async getRating() {
    return { rating: 2 };
  }

  async setRating() {
    return null;
  }

  async getRefreshToken() {
    const data = `${this.dataDir}/refresh_token.json`;

    const ret = () => JSON.parse(fs.readFileSync(data).toString());
    return new Promise<ITokenInfo>((resolve) => {
      setTimeout(() => resolve(ret()), 1000);
    });
  }

  async getAccessToken() {
    const data = `${this.dataDir}/access_token.json`;

    const ret = () => JSON.parse(fs.readFileSync(data).toString());
    return new Promise<ITokenInfo>((resolve) => {
      setTimeout(() => resolve(ret()), 1000);
    });
  }

  async profile() {
    const data = `${this.dataDir}/userInfo.json`;

    const ret = () => JSON.parse(fs.readFileSync(data).toString());
    return new Promise<IUserProfile>((resolve) => {
      setTimeout(() => resolve(ret()), 1000);
    });
  }

  async providers() {
    const data = `${this.dataDir}/provider.json`;
    const ret = () => JSON.parse(fs.readFileSync(data).toString());
    return new Promise<IProvider>((resolve) => {
      setTimeout(() => resolve(ret()), 1000);
    });
  }
}
