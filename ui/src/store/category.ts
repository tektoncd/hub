import { types, getEnv, flow, Instance } from 'mobx-state-tree';
import { Api } from '../api';

export const Tag = types.model('Tags', {
  id: types.identifierNumber,
  name: types.string
});

export const Category = types
  .model('Category', {
    id: types.identifierNumber,
    name: types.string,
    tags: types.array(types.reference(Tag)),
    selected: false
  })

  .actions((self) => ({
    toggle() {
      self.selected = !self.selected;
    }
  }));

export type ITag = Instance<typeof Tag>;
export type ICategory = Instance<typeof Category>;
export type ICategoryStore = Instance<typeof CategoryStore>;

export const CategoryStore = types
  .model('CategoryStore', {
    items: types.map(Category),
    tags: types.optional(types.map(Tag), {}),
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

    get selectedTags() {
      return new Set(
        Array.from(self.items.values())
          .filter((c: ICategory) => c.selected)
          .reduce((acc: number[], c: ICategory) => [...acc, ...c.tags.map((t: ITag) => t.id)], [])
      );
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
    }
  }))

  .actions((self) => ({
    load: flow(function* () {
      try {
        self.setLoading(true);
        const { api } = self;

        const json = yield api.categories();

        // adding the tags to the store - normalized
        const tags: ITag[] = json.data.flatMap((item: ICategory) => item.tags);

        tags.forEach((t) => self.tags.put(t));

        // creating the model only after the store has the tags normalized
        const categories: ICategory[] = json.data.map((c: ICategory) => ({
          id: c.id,
          name: c.name,
          tags: c.tags.map((tag: ITag) => tag.id)
        }));

        categories.forEach((c: ICategory) => self.add(c));
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
