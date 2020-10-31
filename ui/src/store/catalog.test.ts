import { Catalog, CatalogStore } from './catalog';
import { getSnapshot } from 'mobx-state-tree';
import { assert } from './utils';

describe('Store Object', () => {
  it('can create a catalog object', () => {
    const store = Catalog.create({
      id: 1,
      name: 'tekton',
      type: 'community'
    });

    expect(store.name).toBe('tekton');
  });

  it('creates a catalog store', (done) => {
    const store = CatalogStore.create({});

    const item = Catalog.create({
      id: 1,
      name: 'tekton',
      type: 'community'
    });

    store.add(item);

    expect(getSnapshot(store.items)).toMatchSnapshot();

    done();
  });

  it('should toggle a selected catalog', (done) => {
    const store = CatalogStore.create({});

    const item = Catalog.create({
      id: 1,
      name: 'tekton',
      type: 'community'
    });

    store.add(item);

    const catalogs = store.items.get('1');
    assert(catalogs);
    catalogs.toggle();

    expect(store.selected.size).toBe(1);
    expect(catalogs.selected).toBe(true);

    done();
  });

  it('should clear all the selected catalog', (done) => {
    const store = CatalogStore.create({});

    const item = Catalog.create({
      id: 1,
      name: 'tekton',
      type: 'community'
    });

    store.add(item);

    const catalogs = store.items.get('1');
    assert(catalogs);
    catalogs.toggle();

    store.clearSelected();

    expect(catalogs.selected).toBe(false);

    done();
  });
});
