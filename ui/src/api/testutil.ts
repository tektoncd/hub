import { Api } from './';
import * as fs from 'fs';
import { ICategory } from '../store/category';

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
}
