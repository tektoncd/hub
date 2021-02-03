interface API_CONFIG {
  API_URL: string;
  GH_CLIENT_ID: string;
}

export declare global {
  interface Window {
    config: API_CONFIG;
  }
}
