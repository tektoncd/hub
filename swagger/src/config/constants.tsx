interface API_CONFIG {
  API_URL: string;
}

declare global {
  interface Window {
    config: API_CONFIG;
  }
}

window.config = window.config || {
  API_URL: 'no API_URL  set'
};

export const API_URL = window.config.API_URL + '/schema/swagger.json';
