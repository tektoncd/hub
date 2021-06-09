interface API_CONFIG {
  API_URL: string;
  GH_CLIENT_ID: string;
  API_VERSION: string;
  GHE_URL: string;
}

export declare global {
  interface Window {
    config: API_CONFIG;
  }
}
