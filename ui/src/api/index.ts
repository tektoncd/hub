import axios from 'axios';
import { API_URL } from '../config/constants';
import { ICategory } from '../store/category';

export interface Api {
  categories(): Promise<ICategory>;
}

export class Hub implements Api {
  async categories() {
    try {
      return axios.get(`${API_URL}/categories`);
    } catch (err) {
      return err.response;
    }
  }
}
