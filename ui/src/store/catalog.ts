import { Instance, types } from 'mobx-state-tree';
import { Icons } from '../common/icons';

const icons: { [catalog: string]: Icons } = {
  official: Icons.Cat,
  verified: Icons.Certificate,
  community: Icons.User
};

export const Catalog = types
  .model({
    id: types.identifierNumber,
    name: types.optional(types.string, ''),
    type: types.optional(types.string, ''),
    selected: false
  })
  .actions((self) => ({
    toggle() {
      self.selected = !self.selected;
    }
  }))
  .views((self) => ({
    get icon(): Icons {
      return icons[self.type] || Icons.Unknown;
    }
  }));

export type ICatalog = Instance<typeof Catalog>;
export type ICatalogStore = Instance<typeof CatalogStore>;

export const CatalogStore = types
  .model({
    items: types.map(Catalog)
  })

  .actions((self) => ({
    add(item: ICatalog) {
      self.items.put({ id: item.id, name: item.name, type: item.type });
    },
    clearSelected() {
      self.items.forEach((c) => {
        c.selected = false;
      });
    }
  }))

  .views((self) => ({
    get values() {
      return Array.from(self.items.values());
    },

    get selected() {
      const list = new Set();
      self.items.forEach((c: ICatalog) => {
        if (c.selected) {
          list.add(c.id);
        }
      });

      return list;
    }
  }));
