import { TrimUrl } from '../common/trimUrl';

window.config = window.config || {
  API_URL: 'no API_URL  set',
  API_VERSION: 'no API_VERSION set',
  AUTH_BASE_URL: 'no AUTH_BASE_URL set',
  REDIRECT_URI: 'no REDIRECT_URI set'
};

export const API_URL = TrimUrl(window.config.API_URL);
export const API_VERSION = window.config.API_VERSION;
export const AUTH_BASE_URL = TrimUrl(window.config.AUTH_BASE_URL);
export const REDIRECT_URI = TrimUrl(window.config.REDIRECT_URI);
