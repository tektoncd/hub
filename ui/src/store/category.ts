import { AxiosError } from 'axios';
import { types, getEnv, flow, Instance } from 'mobx-state-tree';
import { Api } from '../api';

export const Category = types
  .model('Category', {
    id: types.identifierNumber,
    name: types.string,
    selected: false
  })

  .actions((self) => ({
    toggle() {
      self.selected = !self.selected;
    }
  }));

export type ICategory = Instance<typeof Category>;
export type ICategoryStore = Instance<typeof CategoryStore>;

export const CategoryStore = types
  .model('CategoryStore', {
    items: types.map(Category),
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
    },

    // This returns list of selected category's name
    get selectedByName() {
      return Array.from(self.items.values())
        .filter((c: ICategory) => c.selected)
        .reduce((acc: string[], c: ICategory) => [...acc, c.name], []);
    },

    // This returns set of selected category
    get selected(): Set<number> {
      const list: Set<number> = new Set();
      self.items.forEach((c: ICategory) => {
        if (c.selected) {
          list.add(c.id);
        }
      });
      return list;
    }
  }))

  .actions((self) => ({
    add(item: ICategory) {
      self.items.put(item);
    },

    setLoading(l: boolean) {
      self.isLoading = l;
    },

    clearSelected() {
      self.items.forEach((c) => {
        c.selected = false;
      });
    },

    toggleByName(name: string) {
      self.items.forEach((c) => {
        if (c.name === name) {
          c.selected = true;
        }
      });
    }
  }))

  .actions((self) => ({
    load: flow(function* () {
      try {
        self.setLoading(true);
        const { api } = self;

        const json = yield api.categories();

        // creating the model only after the store has the tags normalized
        const categories: ICategory[] = json.data.map((c: ICategory) => ({
          id: c.id,
          name: c.name
        }));

        categories.forEach((c: ICategory) => self.add(c));
      } catch (error) {
        const err = error as AxiosError;
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
