import { types, Instance } from 'mobx-state-tree';
import { cast } from 'mobx-state-tree';
import fuzzysort from 'fuzzysort';
import { AxiosError } from 'axios';
import moment, { Moment } from 'moment';
import { flow, getEnv } from 'mobx-state-tree';
import { ICategoryStore, Category, ICategory } from './category';
import { Api } from '../api';
import { Catalog, ICatalogStore, ICatalog } from './catalog';
import { Kind, KindStore } from './kind';
import { Platform, PlatformStore, IPlatform } from './platform';
import { assert } from './utils';
import { Params } from '../common/params';
import { Tag, TagStore, ITag } from './tag';
import { TagsKeyword } from '../containers/Search';
import { apiDownError, serverError, resourceNotFoundError } from '../common/errors';

export const updatedAt = types.custom<string, Moment>({
  name: 'momentDate',
  fromSnapshot(value: string): Moment {
    return moment(new Date(value));
  },
  toSnapshot(value: Moment): string {
    return value.fromNow();
  },
  isTargetType(v: string) {
    return moment.isMoment(v);
  },
  getValidationMessage(v: string) {
    if (moment.isMoment(v)) {
      return 'Invalid moment object';
    }
    return '';
  }
});

const Version = types.model('Version', {
  id: types.identifierNumber,
  version: types.string,
  displayName: types.optional(types.string, ''),
  versionPlatforms: types.array(types.reference(Platform)),
  deprecated: types.optional(types.boolean, false),
  description: types.optional(types.string, ''),
  minPipelinesVersion: types.optional(types.string, ''),
  rawURL: types.string,
  webURL: types.string,
  updatedAt: types.optional(updatedAt, '')
});

export const Resource = types
  .model('Resource', {
    id: types.number,
    name: types.string,
    resourceKey: types.identifier,
    catalog: types.reference(Catalog),
    kind: types.reference(Kind),
    categories: types.array(types.reference(Category)),
    platforms: types.array(types.reference(Platform)),
    latestVersion: types.reference(Version),
    displayVersion: types.reference(Version),
    tags: types.array(types.reference(Tag)), // ["cli", "aws"]
    rating: types.number,
    versions: types.array(types.reference(Version)),
    displayName: '',
    tagsString: '',
    readme: '',
    yaml: '',
    status: 0
  })
  .views((self) => ({
    get resourceName() {
      if (self.displayName !== '') {
        return self.displayName;
      }
      return self.name;
    },
    get webURL() {
      const index = self.displayVersion.webURL.lastIndexOf('/');
      return self.displayVersion.webURL.substring(0, index + 1);
    },
    get summary() {
      const description = self.displayVersion.description;
      const index = description.indexOf('\n');
      return description.substring(0, index) || description;
    },
    get detailDescription() {
      const description = self.displayVersion.description;
      const index = description.indexOf('\n');
      return index !== -1 ? description.substring(index + 1).trim() : '';
    },
    get installCommand() {
      return `kubectl apply -f ${self.displayVersion.rawURL}`;
    },
    get tknInstallCommand() {
      const versionFlag =
        self.latestVersion.version !== self.displayVersion.version
          ? ` --version ${self.displayVersion.version}`
          : ``;
      const catalogFlag =
        self.catalog.name.toLowerCase() !== 'tekton'
          ? ` --from ${self.catalog.name.toLowerCase()}`
          : ``;
      return `tkn hub install ${self.kind.name.toLowerCase()} ${
        self.name
      }${versionFlag}${catalogFlag}`;
    }
  }));

export type IResource = Instance<typeof Resource>;
export type IVersion = Instance<typeof Version>;

export enum SortByFields {
  Unknown = '',
  Name = 'Name',
  Rating = 'Rating'
}

export const ResourceStore = types
  .model('ResourceStore', {
    resources: types.map(Resource),
    versions: types.map(Version),
    catalog: types.optional(types.map(Catalog), {}),
    kinds: types.optional(KindStore, {}),
    tags: types.optional(TagStore, {}),
    sortBy: types.optional(types.enumeration(Object.values(SortByFields)), SortByFields.Unknown),
    category: types.optional(types.map(Category), {}),
    platforms: types.optional(PlatformStore, {}),
    search: '',
    searchedTags: types.array(types.string),
    tagsString: '',
    urlParams: '',
    err: '',
    isLoading: true,
    isVersionLoading: true,
    isResourceLoading: true,
    status: 0
  })
  .views((self) => ({
    get items() {
      return Array.from(self.resources.values());
    },
    get api(): Api {
      return getEnv(self).api;
    },
    get categories(): ICategoryStore {
      return getEnv(self).categories;
    },
    get catalogs(): ICatalogStore {
      return getEnv(self).catalogs;
    }
  }))

  .actions((self) => ({
    setLoading(l: boolean) {
      self.isLoading = l;
    },
    setVersionLoading(l: boolean) {
      self.isVersionLoading = l;
    },
    setResourceLoading(l: boolean) {
      self.isResourceLoading = l;
    },
    setSearch(text: string) {
      self.search = text;
    },
    add(item: IResource) {
      self.resources.put(item);
    },
    setSortBy(field: string) {
      const key: SortByFields = SortByFields[field as keyof typeof SortByFields];
      self.sortBy = key;
    },
    setURLParams(url: string) {
      self.urlParams = url;
    },
    setSearchedTags(tags: Array<string>) {
      self.searchedTags = cast(tags);
    }
  }))

  .actions((self) => ({
    clearAllFilters() {
      self.kinds.clearSelected();
      self.catalogs.clearSelected();
      self.categories.clearSelected();
      self.platforms.clearSelected();
      self.setSearch('');
      self.setSearchedTags([]);
      self.setSortBy(SortByFields.Unknown);
    },

    parseUrl() {
      const searchParams = new URLSearchParams(self.urlParams);
      if (searchParams.has(Params.Category)) {
        const categoriesParams = searchParams.getAll(Params.Category)[0].split(',');
        categoriesParams.forEach((t: string) => {
          self.categories.toggleByName(t);
        });
      }

      if (searchParams.has(Params.Catalog)) {
        const catalogsParams = searchParams.getAll(Params.Catalog)[0].split(',');
        catalogsParams.forEach((t: string) => {
          self.catalogs.toggleByName(t);
        });
      }

      if (searchParams.has(Params.Platform)) {
        const platformsParams = searchParams.getAll(Params.Platform)[0].split(',');
        platformsParams.forEach((t: string) => {
          self.platforms.toggleByName(t);
        });
      }

      if (searchParams.has(Params.Kind)) {
        const kindsParams = searchParams.getAll(Params.Kind)[0].split(',');
        kindsParams.forEach((t: string) => {
          const kind = self.kinds.items.get(t);
          assert(kind);
          kind.toggle();
        });
      }
    }
  }))

  .actions((self) => ({
    versionInfo: flow(function* (resourceKey: string) {
      try {
        const { api } = self;
        const resource = self.resources.get(resourceKey);
        assert(resource);

        const json = yield api.resourceVersion(resource.id);

        const versions: IVersion[] = json.data.versions.map((v: IVersion) => ({
          id: v.id,
          version: v.version,
          webURL: v.webURL,
          rawURL: v.rawURL
        }));

        versions.forEach((v: IVersion) => {
          if (!self.versions.has(String(v.id))) {
            self.versions.put(v);
            if (self.resources.has(resourceKey)) {
              const resource = self.resources.get(resourceKey);
              assert(resource);
              resource.versions.push(v.id);
            }
          }
        });
      } catch (error) {
        const err = error as AxiosError;
        self.err = err.toString();
      }
      self.setVersionLoading(false);
    }),

    versionUpdate: flow(function* (versionId: number) {
      try {
        self.setLoading(true);

        const { api } = self;
        const json = yield api.versionUpdate(versionId);

        const versionData = json.data;

        const version: IVersion = {
          id: versionData.id,
          version: versionData.version,
          displayName: versionData.displayName,
          versionPlatforms:
            versionData.platforms != null ? versionData.platforms.map((p: IPlatform) => p.id) : [],
          deprecated: versionData.deprecated || false,
          description: versionData.description,
          minPipelinesVersion: versionData.minPipelinesVersion,
          webURL: versionData.webURL,
          rawURL: versionData.rawURL,
          updatedAt: versionData.updatedAt
        };

        self.versions.put(version);
      } catch (error) {
        const err = error as AxiosError;
        self.err = err.toString();
      }
      self.setLoading(false);
    }),

    load: flow(function* () {
      try {
        self.setLoading(true);

        const { api } = self;

        const json = yield api.resources();
        switch (true) {
          case json.data === undefined:
            self.status = 503;
            self.err = apiDownError;
            break;
          case json.status !== undefined:
            switch (json.status) {
              case 404:
                self.status = json.status;
                self.err = resourceNotFoundError;
                break;
              case 500:
                self.status = json.status;
                self.err = serverError;
                break;
            }
            break;
          default:
            {
              const kinds: string[] = json.data.map((r: IResource) => r.kind);
              kinds.forEach((k) => self.kinds.add(k));

              json.data.forEach((r: IResource) => {
                self.versions.put(r.latestVersion);
              });
              // Adding the tags to the store - normalized
              const tags: ITag[] = json.data.flatMap((item: IResource) => item.tags);
              tags.forEach((t) => self.tags.add(t));

              const allCatalogs: ICatalog[] = json.data.flatMap((item: IResource) => item.catalog);
              allCatalogs.forEach((t) => self.catalog.put(t));

              //Adding the categories to the store - normalized
              const categories: ICategory[] = json.data.flatMap(
                (item: IResource) => item.categories
              );
              categories.forEach((c) => {
                self.category.put(c);
              });

              const allPlatforms: IPlatform[] = json.data.flatMap(
                (item: IResource) => item.platforms
              );
              allPlatforms.forEach((p) => self.platforms.add(p));

              const resources: IResource[] = json.data.map((r: IResource) => ({
                id: r.id,
                name: r.name,
                resourceKey: `${r.catalog.name}/${r.kind}/${r.name}`,
                catalog: r.catalog.id,
                kind: r.kind,
                latestVersion: r.latestVersion.id,
                displayVersion: r.latestVersion.id,
                tags: r.tags != null ? r.tags.map((tag: ITag) => tag.name) : [],
                categories: r.categories != null ? r.categories.map((c: ICategory) => c.id) : [],
                platforms: r.platforms != null ? r.platforms.map((p: IPlatform) => p.id) : [],
                rating: r.rating,
                versions: [],
                displayName: r.latestVersion.displayName,
                tagsString: r.tags.map((tag: ITag) => tag.name).join(' ')
              }));

              resources.forEach((r: IResource) => {
                r.versions.push(r.latestVersion);
                self.add(r);
              });
            }
            break;
        }
        // Url parsing after resource load
        if (self.urlParams) self.parseUrl();
      } catch (error) {
        const err = error as AxiosError;
        self.err = err.toString();
      }
      self.setLoading(false);
      self.setResourceLoading(false);
    }),

    loadReadme: flow(function* (name: string) {
      try {
        self.setLoading(true);

        const { api, resources } = self;
        const resource = resources.get(name);
        assert(resource);
        const version = resource.displayVersion.version;
        assert(version);

        const readme = yield api.readme(name, version);
        resource.readme = readme;
      } catch (error) {
        const err = error as AxiosError;
        self.err = err.toString();
      }
      self.setLoading(false);
    }),

    loadYaml: flow(function* (name: string) {
      try {
        self.setLoading(true);

        const { api, resources } = self;
        const resource = resources.get(name);
        assert(resource);
        const version = resource.displayVersion.version;
        assert(version);

        const yaml = yield api.yaml(name, version);
        resource.yaml = yaml;
      } catch (error) {
        const err = error as AxiosError;
        self.err = err.toString();
      }
      self.setLoading(false);
    })
  }))

  .actions((self) => ({
    afterCreate() {
      self.load();
    },
    setDisplayVersion(resourceKey: string, versionId: string | number) {
      const resource = self.resources.get(resourceKey);
      assert(resource);
      const version = self.versions.get(versionId as string);
      assert(version);
      if (version.id !== resource.displayVersion.id) {
        resource.displayVersion = version;
        self.versionUpdate(version.id);
        self.loadReadme(resourceKey);
        self.loadYaml(resourceKey);
      }
    }
  }))

  .views((self) => ({
    get filteredResources() {
      const { resources, kinds, catalogs, platforms, categories } = self;
      const { search, sortBy, searchedTags } = self;

      const tags = new Set(searchedTags);

      let filteredItems: IResource[] = [];
      resources.forEach((r: IResource) => {
        const matchesKind = kinds.selected.size === 0 || kinds.selected.has(r.kind.name);
        const matchesCatalogs = catalogs.selected.size === 0 || catalogs.selected.has(r.catalog.id);
        const matchesCategories =
          categories.selected.size === 0 || r.categories.some((c) => categories.selected.has(c.id));
        const matchesPlatforms =
          platforms.selected.size === 0 || r.platforms.some((p) => platforms.selected.has(p.id));
        const matchesTags = searchedTags.length === 0 || r.tags.some((t: ITag) => tags.has(t.name));

        if (
          matchesKind &&
          matchesCatalogs &&
          matchesCategories &&
          matchesTags &&
          matchesPlatforms
        ) {
          filteredItems.push(r);
        }
      });

      if (search == TagsKeyword) {
        return filteredItems;
      }

      if (search.trim() !== '' && searchedTags.length === 0) {
        filteredItems = fuzzysort
          .go(search, filteredItems, { keys: ['name', 'displayName', 'tagsString'] })
          .map((resource: Fuzzysort.KeysResult<IResource>) => resource.obj);
      }

      switch (sortBy) {
        case SortByFields.Rating:
          return filteredItems.sort((first: IResource, second: IResource) =>
            first.rating < second.rating ? 1 : first.rating > second.rating ? -1 : 0
          );

        case SortByFields.Name:
          return filteredItems.sort((first: IResource, second: IResource) =>
            first.name > second.name ? 1 : first.name < second.name ? -1 : 0
          );

        default:
          return filteredItems;
      }
    }
  }));

export type IResourceStore = Instance<typeof ResourceStore>;
