import { ResourceStore, Resource, SortByFields } from './resource';
import { getSnapshot } from 'mobx-state-tree';
import { when } from 'mobx';
import { FakeHub } from '../api/testutil';
import { CategoryStore } from './category';
import { assert } from './utils';

const TESTDATA_DIR = `${__dirname}/testdata`;
const api = new FakeHub(TESTDATA_DIR);

describe('Store Object', () => {
  it('can create a resource object', () => {
    const store = Resource.create({
      id: 5,
      name: 'buildah',
      catalog: '1',
      kind: 'Task',
      latestVersion: 1,
      tags: ['1'],
      rating: 5
    });

    expect(store.name).toBe('buildah');
  });
});

describe('Store functions', () => {
  it('creates a resource store', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        categories: CategoryStore.create({}, { api })
      }
    );

    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.resources.size).toBe(6);
        expect(getSnapshot(store.resources)).toMatchSnapshot();
        done();
      }
    );
  });

  it('creates a catalog store', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.resources.size).toBe(6);
        expect(getSnapshot(store.catalogs)).toMatchSnapshot();
        done();
      }
    );
  });

  it('creates a kind store', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        categories: CategoryStore.create({}, { api })
      }
    );

    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.resources.size).toBe(6);
        expect(getSnapshot(store.kinds)).toMatchSnapshot();
        done();
      }
    );
  });

  it('filter resources based on selected catalog', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.resources.size).toBe(6);
        const { items } = store.catalogs;

        const catalogs = items.get('2');
        assert(catalogs);
        catalogs.toggle();

        const filtered = store.filteredResources;
        expect(filtered.length).toBe(1);
        expect(filtered[0].name).toBe('hub');
        expect(filtered[0].catalog.name).toBe('tekton-hub');

        done();
      }
    );
  });

  it('filter resources based on selected kind', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.resources.size).toBe(6);

        const kind = store.kinds.items.get('Pipeline');
        assert(kind);
        kind.toggle();

        expect(store.filteredResources.length).toBe(1);
        expect(store.filteredResources[0].name).toBe('hub');
        expect(store.filteredResources[0].kind.name).toBe('Pipeline');

        done();
      }
    );
  });

  it('filter resources based on selected kind and catalog', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.resources.size).toBe(6);

        const kinds = store.kinds.items.get('Task');
        assert(kinds);

        const catalogs = store.catalogs.items.get('1');
        assert(catalogs);

        kinds.toggle();
        catalogs.toggle();

        store.setSearch('golang');

        store.filteredResources;

        expect(store.filteredResources.length).toBe(1);

        done();
      }
    );
  });

  it('makes sure to not add duplicate resources', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        categories: CategoryStore.create({}, { api })
      }
    );

    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        const item = Resource.create({
          id: 44,
          name: 'golang-build',
          catalog: 1,
          kind: 'Task',
          latestVersion: 47,
          tags: [1],
          rating: 5,
          versions: [47],
          displayName: 'golang build'
        });

        store.add(item);
        expect(store.resources.size).toBe(6);

        expect(getSnapshot(store.resources)).toMatchSnapshot();

        done();
      }
    );
  });

  it('it checks if the related date is a string', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        expect(store.isLoading).toBe(false);
        expect(store.resources.size).toBe(6);

        const versions = store.versions.get('1');
        assert(versions);
        expect(typeof versions.updatedAt.fromNow()).toBe('string');

        done();
      }
    );
  });

  it('should filter resources based on search', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        store.setSearch('golang');
        expect(store.filteredResources[0].name).toBe('golang-build');

        done();
      }
    );
  });

  it('it should return displayName', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        expect(store.isLoading).toBe(false);
        expect(store.resources.size).toBe(6);

        const displayName = store.resources.get('aws-cli')?.resourceName;

        expect(displayName).toBe('aws cli');
        done();
      }
    );
  });

  it('should sort resources based on selected key', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        const rating: SortByFields = SortByFields[SortByFields.Rating];
        store.setSortBy(rating);

        expect(store.filteredResources[0].rating).toBe(5);
        expect(store.filteredResources[0].name).toBe('aws-cli');

        const name: SortByFields = SortByFields[SortByFields.Name];
        store.setSortBy(name);

        expect(store.filteredResources[0].name).toBe('ansible-runner');

        done();
      }
    );
  });
});
