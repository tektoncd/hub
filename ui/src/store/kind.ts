import { types, Instance } from 'mobx-state-tree';

export const Kind = types
  .model({
    name: types.identifier,
    selected: false
  })
  .actions((self) => ({
    toggle() {
      self.selected = !self.selected;
    }
  }));

export type IKind = Instance<typeof Kind>;
export type IKindStore = Instance<typeof KindStore>;

export const KindStore = types
  .model({
    items: types.optional(types.map(Kind), {})
  })

  .actions((self) => ({
    add(item: string) {
      self.items.put({ name: item, selected: false });
    },

    clearSelected() {
      self.items.forEach((k) => {
        k.selected = false;
      });
    }
  }))

  .views((self) => ({
    get values() {
      return Array.from(self.items.values());
    },

    get selected() {
      const list = new Set();
      self.items.forEach((c: IKind) => {
        if (c.selected) {
          list.add(c.name);
        }
      });

      return list;
    }
  }));
