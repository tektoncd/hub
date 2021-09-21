import { types, getEnv, flow, Instance } from 'mobx-state-tree';
import { Api } from '../api';

export const Provider = types.model('Provider', {
  name: types.identifier
});

export type IProvider = Instance<typeof Provider>;
export type IProviderStore = Instance<typeof ProviderStore>;

export const ProviderStore = types
  .model('ProviderStore', {
    items: types.map(Provider),
    isLoading: true,
    err: ''
  })

  .views((self) => ({
    get api(): Api {
      return getEnv(self).api;
    },

    get count() {
      return self.items.size;
    },

    get values() {
      return Array.from(self.items.values());
    }
  }))

  .actions((self) => ({
    add(item: IProvider) {
      self.items.put(item);
    },

    setLoading(l: boolean) {
      self.isLoading = l;
    }
  }))

  .actions((self) => ({
    load: flow(function* () {
      try {
        self.setLoading(true);
        const { api } = self;

        const json = yield api.providers();

        const providrs: IProvider[] = json.data.map((c: IProvider) => ({
          name: c.name
        }));

        providrs.forEach((c: IProvider) => self.add(c));
      } catch (err) {
        self.err = err.toString();
      }

      self.setLoading(false);
    })
  }))

  .actions((self) => ({
    afterCreate() {
      self.load();
    }
  }));
