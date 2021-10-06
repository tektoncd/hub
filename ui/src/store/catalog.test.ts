import { getSnapshot } from 'mobx-state-tree';
import { when } from 'mobx';
import { Catalog, CatalogStore } from './catalog';
import { assert } from './utils';
import { Icons } from '../common/icons';
import { FakeHub } from '../api/testutil';

const TESTDATA_DIR = `${__dirname}/testdata`;
const api = new FakeHub(TESTDATA_DIR);

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
    const store = CatalogStore.create({}, { api });
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        expect(store.items.size).toBe(2);
        expect(store.isLoading).toBe(false);

        expect(getSnapshot(store)).toMatchSnapshot();

        done();
      }
    );
  });

  it('should toggle a selected catalog', (done) => {
    const store = CatalogStore.create({}, { api });
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        expect(store.isLoading).toBe(false);
        const catalogs = store.items.get('1');
        assert(catalogs);
        catalogs.toggle();

        expect(store.selected.size).toBe(1);
        expect(catalogs.selected).toBe(true);

        done();
      }
    );
  });

  it('should clear all the selected catalog', (done) => {
    const store = CatalogStore.create({}, { api });
    when(
      () => !store.isLoading,
      () => {
        expect(store.isLoading).toBe(false);

        const catalogs = store.items.get('1');
        assert(catalogs);
        catalogs.toggle();

        store.clearSelected();

        expect(catalogs.selected).toBe(false);
        done();
      }
    );
  });

  it('can get an icon for catalog', (done) => {
    const store = CatalogStore.create({}, { api });

    when(
      () => !store.isLoading,
      () => {
        const catalog = store.items.get('1');
        assert(catalog);
        expect(catalog.icon).toBe(Icons.Catalog);

        done();
      }
    );
  });

  it('should toggle catalogs by name and can get selected catlogs by name', (done) => {
    const store = CatalogStore.create({}, { api });

    when(
      () => !store.isLoading,
      () => {
        store.toggleByName('Tekton');
        const catalogs = store.items.get('1');
        assert(catalogs);

        expect(catalogs.selected).toBe(true);
        expect(store.selectedByName).toEqual(['tekton']);

        done();
      }
    );
  });

  it('can get a provider for catalog', (done) => {
    const store = CatalogStore.create({}, { api });

    when(
      () => !store.isLoading,
      () => {
        const catalog = store.items.get('2');
        assert(catalog);
        expect(catalog.provider).toBe('gitlab');

        done();
      }
    );
  });
});
