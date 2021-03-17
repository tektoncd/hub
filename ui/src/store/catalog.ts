import { Instance, types } from 'mobx-state-tree';
import { Icons } from '../common/icons';
import { titleCase } from '../common/titlecase';

const icons: { [catalog: string]: Icons } = {
  community: Icons.Catalog
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
      return icons[self.type] || Icons.Catalog;
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
    },
    toggleByName(name: string) {
      self.items.forEach((c) => {
        if (titleCase(c.name) === name) {
          c.selected = true;
        }
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
    },

    /* This view returns list of the selected catalos's name instead of id
    to avoid loop on it inorder to get catalogs name */
    get selectedByName() {
      return Array.from(self.items.values())
        .filter((c: ICatalog) => c.selected)
        .reduce((acc: string[], c: ICatalog) => [...acc, c.name], []);
    }
  }));
