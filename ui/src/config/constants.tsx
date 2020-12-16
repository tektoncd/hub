interface API_CONFIG {
  API_URL: string;
  GH_CLIENT_ID: string;
}

declare global {
  interface Window {
    config: API_CONFIG;
  }
}

window.config = window.config || {
  API_URL: 'no API_URL  set',
  GH_CLIENT_ID: 'no GH_CLIENT_ID set'
};

export const API_URL = window.config.API_URL;
export const GH_CLIENT_ID = window.config.GH_CLIENT_ID;
