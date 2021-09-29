import { Instance, types, flow, getEnv } from 'mobx-state-tree';
import { Icons } from '../common/icons';
import { titleCase } from '../common/titlecase';
import { Api } from '../api';

const icons: { [catalog: string]: Icons } = {
  community: Icons.Catalog
};

export const Catalog = types
  .model({
    id: types.identifierNumber,
    name: types.optional(types.string, ''),
    type: types.optional(types.string, ''),
    provider: types.optional(types.string, ''),
    selected: false
  })
  .actions((self) => ({
    toggle() {
      self.selected = !self.selected;
    }
  }))
  .views((self) => ({
    get icon(): Icons {
      return icons[self.type] || Icons.Catalog;
    }
  }));

export type ICatalog = Instance<typeof Catalog>;
export type ICatalogStore = Instance<typeof CatalogStore>;

export const CatalogStore = types
  .model({
    items: types.map(Catalog),
    isLoading: true,
    err: ''
  })

  .actions((self) => ({
    add(item: ICatalog) {
      self.items.put({ id: item.id, name: item.name, type: item.type, provider: item.provider });
    },
    clearSelected() {
      self.items.forEach((c) => {
        c.selected = false;
      });
    },
    toggleByName(name: string) {
      self.items.forEach((c) => {
        if (titleCase(c.name) === name) {
          c.selected = true;
        }
      });
    },
    setLoading(l: boolean) {
      self.isLoading = l;
    }
  }))

  .views((self) => ({
    get values() {
      return Array.from(self.items.values());
    },

    get api(): Api {
      return getEnv(self).api;
    },

    get selected() {
      const list = new Set();
      self.items.forEach((c: ICatalog) => {
        if (c.selected) {
          list.add(c.id);
        }
      });

      return list;
    },

    /* This view returns list of the selected catalos's name instead of id
    to avoid loop on it inorder to get catalogs name */
    get selectedByName() {
      return Array.from(self.items.values())
        .filter((c: ICatalog) => c.selected)
        .reduce((acc: string[], c: ICatalog) => [...acc, c.name], []);
    }
  }))

  .actions((self) => ({
    load: flow(function* () {
      try {
        self.setLoading(true);
        const { api } = self;

        const json = yield api.catalogs();

        const catalogs: ICatalog[] = json.data.map((c: ICatalog) => ({
          id: c.id,
          name: c.name,
          type: c.type,
          provider: c.provider
        }));

        catalogs.forEach((c: ICatalog) => self.add(c));
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
