import { ResourceStore, Resource, SortByFields } from './resource';
import { getSnapshot } from 'mobx-state-tree';
import { when } from 'mobx';
import { FakeHub } from '../api/testutil';
import { CategoryStore } from './category';
import { CatalogStore } from './catalog';
import { assert } from './utils';

const TESTDATA_DIR = `${__dirname}/testdata`;
const api = new FakeHub(TESTDATA_DIR);

describe('Store Object', () => {
  it('can create a resource object', () => {
    const store = Resource.create({
      id: 5,
      name: 'buildah',
      resourceKey: 'tekton/Task/buildah',
      catalog: '1',
      kind: 'Task',
      latestVersion: 1,
      displayVersion: 1,
      tags: ['1'],
      platforms: ['1'],
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
        expect(store.resources.size).toBe(7);
        expect(getSnapshot(store.resources)).toMatchSnapshot();
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
        expect(store.resources.size).toBe(7);
        expect(getSnapshot(store.kinds)).toMatchSnapshot();
        done();
      }
    );
  });

  it('creates a platform store', (done) => {
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
        expect(store.resources.size).toBe(7);
        expect(getSnapshot(store.platforms)).toMatchSnapshot();
        done();
      }
    );
  });

  it('filter resources based on selected catalog', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        catalogs: CatalogStore.create({}, { api }),
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.resources.size).toBe(7);
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
        catalogs: CatalogStore.create({}, { api }),
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.resources.size).toBe(7);

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

  it('filter resources based on selected categories', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        catalogs: CatalogStore.create({}, { api }),
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        const caetgoryA = store.categories.items.get('1');
        const categoryB = store.categories.items.get('2');

        assert(caetgoryA);
        assert(categoryB);

        caetgoryA.toggle();
        categoryB.toggle();

        expect(store.filteredResources.length).toBe(2);
        expect(store.filteredResources[0].name).toBe('golang-build');
        done();
      }
    );
  });

  it('filter resources based on selected platforms', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        catalogs: CatalogStore.create({}, { api }),
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        const platformA = store.platforms.items.get('2');
        const platformB = store.platforms.items.get('3');

        assert(platformA);
        assert(platformB);

        platformA.toggle();
        platformB.toggle();

        expect(store.filteredResources.length).toBe(2);
        expect(store.filteredResources[0].name).toBe('buildah');
        done();
      }
    );
  });

  it('filter resources based on selected kind and catalog', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        catalogs: CatalogStore.create({}, { api }),
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.resources.size).toBe(7);

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
          resourceKey: 'tekton/Task/golang-build',
          catalog: 1,
          kind: 'Task',
          latestVersion: 47,
          displayVersion: 47,
          tags: [1],
          rating: 5,
          versions: [47],
          displayName: 'golang build'
        });

        store.add(item);
        expect(store.resources.size).toBe(7);

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
        expect(store.resources.size).toBe(7);

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
        catalogs: CatalogStore.create({}, { api }),
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        store.setSearch('golang');
        expect(store.filteredResources[0].name).toBe('golang-build');

        store.setSearch('resource');
        expect(store.filteredResources[0].name).toBe('hub');
        expect(store.filteredResources[0].tags[0].name).toBe('resource');

        done();
      }
    );
  });

  it('should filter return all resources if searched text is `tags:`', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        catalogs: CatalogStore.create({}, { api }),
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        store.setSearch('tags:');

        expect(store.filteredResources.length).toBe(7);

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
        expect(store.resources.size).toBe(7);

        const resource = store.resources.get('tekton/Task/aws-cli');
        assert(resource);

        expect(resource.resourceName).toBe('aws cli');
        done();
      }
    );
  });

  it('should sort resources based on selected key', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        catalogs: CatalogStore.create({}, { api }),
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

  it('should sort resources based on last updated at', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        catalogs: CatalogStore.create({}, { api }),
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        const lastUpdated: SortByFields = SortByFields[SortByFields.RecentlyUpdated];
        store.setSortBy(lastUpdated);

        expect(store.filteredResources[0].rating).toBe(4);
        expect(store.filteredResources[0].name).toBe(`buildah`);

        done();
      }
    );
  });

  it('it should return webURL', (done) => {
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
        expect(store.resources.size).toBe(7);

        const resource = store.resources.get('tekton/Task/aws-cli');
        assert(resource);

        expect(resource.webURL).toBe(
          'https://github.com/tektoncd/catalog/tree/master/task/aws-cli/0.1/'
        );
        done();
      }
    );
  });

  it('it should return summary', (done) => {
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
        expect(store.resources.size).toBe(7);

        const resource = store.resources.get('tekton/Task/aws-cli');
        assert(resource);

        expect(resource.summary).toBe(
          'This task performs operations on Amazon Web Services resources using aws.'
        );
        done();
      }
    );
  });

  it('it should return detail description', (done) => {
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
        expect(store.resources.size).toBe(7);

        const resource = store.resources.get('tekton/Task/buildah');
        assert(resource);

        expect(resource.detailDescription).toBe(
          "Buildah Task builds source into a container image using Project Atomic's Buildah build tool.It uses Buildah's support for building from Dockerfiles, using its buildah bud command.This command executes the directives in the Dockerfile to assemble a container image, then pushes that image to a container registry."
        );
        done();
      }
    );
  });

  it('it should return install command', (done) => {
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
        expect(store.resources.size).toBe(7);

        const resource = store.resources.get('tekton/Task/aws-cli');
        assert(resource);

        expect(resource.installCommand).toBe(
          'kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/master/task/aws-cli/0.1/aws-cli.yaml'
        );
        done();
      }
    );
  });

  it('update versions list for buildah resource', (done) => {
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
        expect(store.resources.size).toBe(7);
        expect(getSnapshot(store.resources)).toMatchSnapshot();
        store.versionInfo('tekton/Task/buildah');
        when(
          () => !store.isVersionLoading,
          () => {
            const resource = store.resources.get('tekton/Task/buildah');
            assert(resource);
            expect(resource.versions.length).toBe(2);
            done();
          }
        );
      }
    );
  });

  it('fetch 0.1 version details for buildah resource', (done) => {
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
        expect(store.resources.size).toBe(7);
        expect(getSnapshot(store.resources)).toMatchSnapshot();
        store.versionInfo('tekton/Task/buildah');
        when(
          () => !store.isVersionLoading,
          () => {
            const resource = store.resources.get('tekton/Task/buildah');
            assert(resource);
            expect(resource.versions.length).toBe(2);
            store.versionUpdate(13);
            when(
              () => !store.isLoading,
              () => {
                expect(resource.versions[1].minPipelinesVersion).toBe('0.12.1');
                expect(resource.versions[1].version).toBe('0.1');
                done();
              }
            );
          }
        );
      }
    );
  });

  it('fetch platforms version details for buildah resource', (done) => {
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
        expect(store.resources.size).toBe(7);
        expect(getSnapshot(store.resources)).toMatchSnapshot();
        store.versionInfo('tekton/Task/buildah');
        when(
          () => !store.isVersionLoading,
          () => {
            const resource = store.resources.get('tekton/Task/buildah');
            assert(resource);
            expect(resource.versions.length).toBe(2);
            expect(resource.latestVersion.id).toBe(105);
            store.versionUpdate(105);
            when(
              () => !store.isLoading,
              () => {
                expect(resource.versions[0].versionPlatforms.length).toBe(2);
                expect(resource.versions[0].versionPlatforms[0].name).toBe('linux/amd64');
                expect(resource.versions[0].versionPlatforms[1].name).toBe('linux/s390x');
                expect(resource.versions[1].id).toBe(13);
                store.versionUpdate(13);
                when(
                  () => !store.isLoading,
                  () => {
                    expect(resource.versions[1].versionPlatforms.length).toBe(1);
                    expect(resource.versions[1].versionPlatforms[0].name).toBe('linux/amd64');
                    done();
                  }
                );
              }
            );
          }
        );
      }
    );
  });

  it('set 0.1 as display version for buildah', (done) => {
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
        expect(store.resources.size).toBe(7);
        expect(getSnapshot(store.resources)).toMatchSnapshot();
        store.versionInfo('tekton/Task/buildah');
        when(
          () => !store.isVersionLoading,
          () => {
            const resource = store.resources.get('tekton/Task/buildah');
            assert(resource);
            expect(resource.versions.length).toBe(2);
            store.setDisplayVersion('tekton/Task/buildah', '13');
            when(
              () => !store.isLoading,
              () => {
                expect(resource.versions[1].minPipelinesVersion).toBe('0.12.1');
                expect(resource.versions[1].version).toBe('0.1');
                done();
              }
            );
          }
        );
      }
    );
  });

  it('should display the readme for buildah', (done) => {
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
        expect(store.resources.size).toBe(7);
        store.loadReadme('buildah');
        when(
          () => !store.isLoading,
          () => {
            const resource = store.resources.get('tekton/Task/buildah');
            assert(resource);

            expect(typeof resource.readme).toBe('string');
            done();
          }
        );
      }
    );
  });

  it('should display the yaml for buildah', (done) => {
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
        expect(store.resources.size).toBe(7);
        store.loadYaml('buildah');
        when(
          () => !store.isLoading,
          () => {
            const resource = store.resources.get('tekton/Task/buildah');
            assert(resource);

            expect(typeof resource.readme).toBe('string');
            done();
          }
        );
      }
    );
  });

  it('it should clear all selected filters', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        catalogs: CatalogStore.create({}, { api }),
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        store.setSearch('golang');
        store.setSortBy(SortByFields.Name);
        expect(store.filteredResources.length).toBe(1);

        store.clearAllFilters();
        expect(store.filteredResources.length).toBe(7);

        done();
      }
    );
  });

  it('makes sure to add resources with same name but from different catalog', (done) => {
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
          resourceKey: 'tekton-hub/Task/golang-build',
          catalog: 2,
          kind: 'Task',
          latestVersion: 47,
          displayVersion: 47,
          tags: [1],
          rating: 5,
          versions: [47],
          displayName: 'golang build'
        });

        expect(store.resources.size).toBe(7);
        store.add(item);
        expect(store.resources.size).toBe(8);

        expect(getSnapshot(store.resources)).toMatchSnapshot();

        done();
      }
    );
  });

  it('it should return tkn hub install command', (done) => {
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
        expect(store.resources.size).toBe(7);

        const tektonResource = store.resources.get('tekton/Task/aws-cli');
        assert(tektonResource);

        expect(tektonResource.tknInstallCommand).toBe('tkn hub install task aws-cli');

        const hubResource = store.resources.get('tekton-hub/Pipeline/hub');

        assert(hubResource);
        expect(hubResource.tknInstallCommand).toBe(
          'tkn hub install pipeline hub --from tekton-hub'
        );
        done();
      }
    );
  });

  it('it should parse the url and can update the store', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        catalogs: CatalogStore.create({}, { api }),
        categories: CategoryStore.create({}, { api })
      }
    );
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        store.setURLParams('?category=Automation%2CBuild+Tools&catalog=tekton');
        store.parseUrl();

        expect(store.filteredResources.length).toBe(1);
        expect(store.categories.selectedByName).toEqual(['Build Tools']);

        done();
      }
    );
  });

  it('it should have deprecation info in the displayVersion for a deprecated version', (done) => {
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
        const resource = store.resources.get('tekton/Task/aws-cli');
        assert(resource);

        expect(resource.displayVersion.deprecated).toBe(true);
        done();
      }
    );
  });
});
