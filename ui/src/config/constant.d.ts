interface API_CONFIG {
  API_URL: string;
  API_VERSION: string;
  AUTH_BASE_URL: string;
  REDIRECT_URI: string;
  CUSTOM_LOGO_MEDIA_TYPE: string;
  CUSTOM_LOGO_BASE64_DATA: string;
}

export declare global {
  interface Window {
    config: API_CONFIG;
  }
}
