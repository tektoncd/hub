import { types, Instance } from 'mobx-state-tree';

export const Tag = types.model('Tags', {
  id: types.number,
  name: types.identifier
});

export const TagStore = types
  .model({
    items: types.optional(types.map(Tag), {})
  })
  .actions((self) => ({
    add(item: ITag): void {
      self.items.put(item);
    }
  }))
  .views((self) => ({
    get values() {
      return Array.from(self.items.values());
    }
  }));

export type ITag = Instance<typeof Tag>;
export type ItagStore = Instance<typeof TagStore>;
