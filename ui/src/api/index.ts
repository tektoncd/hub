import axios from 'axios';
import { API_URL } from '../config/constants';
import { ICategory } from '../store/category';
import { IResource } from '../store/resource';

export interface Api {
  categories(): Promise<ICategory>;
  resources(): Promise<IResource>;
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
}
