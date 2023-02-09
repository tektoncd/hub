import React from 'react';
import { types, getEnv, Instance } from 'mobx-state-tree';
import { persist } from 'mst-persist';
import { CategoryStore, ICategoryStore } from './category';
import { ResourceStore, IResourceStore } from './resource';
import { Hub, Api } from '../api';
import { AuthStore, IAuthStore } from './auth';
import { CatalogStore, ICatalogStore } from './catalog';
import { ProviderStore, IProviderStore } from './provider';

export const Root = types.model('Root', {}).views((self) => ({
  get api(): Api {
    return getEnv(self).api;
  },
  get categories(): ICategoryStore {
    return getEnv(self).categories;
  },
  get resources(): IResourceStore {
    return getEnv(self).resources;
  },
  get user(): IAuthStore {
    return getEnv(self).user;
  },
  get catalogs(): ICatalogStore {
    return getEnv(self).catalogs;
  },
  get providers(): IProviderStore {
    return getEnv(self).providers;
  }
}));

export type IRoot = Instance<typeof Root>;

const initRootStore = (api: Api) => {
  const categories = CategoryStore.create({}, { api });
  const catalogs = CatalogStore.create({}, { api });
  const resources = ResourceStore.create({}, { api, categories, catalogs });
  const user = AuthStore.create({ accessTokenInfo: {}, refreshTokenInfo: {} }, { api });
  const providers = ProviderStore.create({}, { api });

  // This peristently stores user auth details
  persist('authStore', user, {
    storage: localStorage,

    // This adds fields from the store for which data needs to be stored persistently
    whitelist: ['accessTokenInfo', 'refreshTokenInfo', 'isAuthenticated', 'isLoading']
  });

  return Root.create({}, { api, catalogs, categories, resources, user, providers });
};

interface Props {
  children: React.ReactNode;
}

let RootContext: React.Context<IRoot>;
export const useMst = () => React.useContext(RootContext);

export const createProviderAndStore = (api?: Api) => {
  const root = initRootStore(api || new Hub());
  RootContext = React.createContext(root);

  const Provider = ({ children }: Props) => (
    <>
      <RootContext.Provider value={root}> {children} </RootContext.Provider>
    </>
  );
  return { Provider, root };
};

export const createProvider = (api?: Api) => {
  const { Provider } = createProviderAndStore(api);
  return Provider;
};
