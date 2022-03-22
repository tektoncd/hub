import { types, Instance } from 'mobx-state-tree';

export const Platform = types
  .model({
    id: types.identifierNumber,
    name: types.string,
    selected: false
  })

  .actions((self) => ({
    toggle() {
      self.selected = !self.selected;
    }
  }));

export type IPlatform = Instance<typeof Platform>;
export type IPlatformStore = Instance<typeof PlatformStore>;

export const PlatformStore = types
  .model({
    items: types.optional(types.map(Platform), {})
  })

  .actions((self) => ({
    add(item: IPlatform): void {
      self.items.put(item);
    },

    clearSelected() {
      self.items.forEach((p) => {
        p.selected = false;
      });
    },

    toggleByName(name: string) {
      self.items.forEach((p) => {
        if (p.name === name) {
          p.selected = true;
        }
      });
    }
  }))

  .views((self) => ({
    get values() {
      return Array.from(self.items.values());
    },

    get selectedByName() {
      return Array.from(self.items.values())
        .filter((p: IPlatform) => p.selected)
        .reduce((acc: string[], p: IPlatform) => [...acc, p.name], []);
    },

    get selected(): Set<number> {
      const list: Set<number> = new Set();
      self.items.forEach((p: IPlatform) => {
        if (p.selected) {
          list.add(p.id);
        }
      });
      return list;
    }
  }));
