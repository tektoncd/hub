import * as fs from 'fs';
import { ICategory } from '../store/category';
import { IResource, IVersion } from '../store/resource';
import { Api, AuthResponse } from './';

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

  async readme(rawURL: string) {
    const splitRawUrl = rawURL.split('/');
    const resourceName = splitRawUrl[splitRawUrl.length - 3];
    const data = `${this.dataDir}/${resourceName}-Readme.md`;

    const ret = () => fs.readFileSync(data).toString();
    return new Promise<string>((resolve) => {
      setTimeout(() => resolve(ret()), 1000);
    });
  }

  async yaml(rawURL: string) {
    const splitRawUrl = rawURL.split('/');
    const resourceName = splitRawUrl[splitRawUrl.length - 3];
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
}
