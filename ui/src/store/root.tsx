import React, { ReactChild, ReactChildren } from 'react';
import { types, getEnv, Instance } from 'mobx-state-tree';
import { persist } from 'mst-persist';
import { CategoryStore } from './category';
import { ResourceStore } from './resource';
import { Hub, Api } from '../api';
import { AuthStore } from './auth';
import { CatalogStore } from './catalog';

export const Root = types.model('Root', {}).views((self) => ({
  get api(): Api {
    return getEnv(self).api;
  },
  get categories() {
    return getEnv(self).categories;
  },
  get resources() {
    return getEnv(self).resources;
  },
  get user() {
    return getEnv(self).user;
  },
  get catalogs() {
    return getEnv(self).catalogs;
  }
}));

type IRoot = Instance<typeof Root>;

const initRootStore = (api: Api) => {
  const categories = CategoryStore.create({}, { api });
  const catalogs = CatalogStore.create({}, { api });
  const resources = ResourceStore.create({}, { api, categories, catalogs });
  const user = AuthStore.create({ accessTokenInfo: {}, refreshTokenInfo: {} }, { api });

  // This peristently stores user auth details
  persist('authStore', user, {
    storage: localStorage,

    // This adds fields from the store for which data needs to be stored persistently
    whitelist: ['accessTokenInfo', 'refreshTokenInfo', 'isAuthenticated', 'isLoading']
  });

  return Root.create({}, { api, catalogs, categories, resources, user });
};

interface Props {
  children: ReactChild | ReactChildren;
}

let RootContext: React.Context<IRoot>;
export const useMst = () => React.useContext(RootContext);

export const createProviderAndStore = (api?: Api) => {
  const root = initRootStore(api || new Hub());
  RootContext = React.createContext(root);

  const Provider = ({ children }: Props) => (
    <RootContext.Provider value={root}> {children} </RootContext.Provider>
  );
  return { Provider, root };
};

export const createProvider = (api?: Api) => {
  const { Provider } = createProviderAndStore(api);
  return Provider;
};
