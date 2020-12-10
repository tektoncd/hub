import { types, Instance } from 'mobx-state-tree';
import fuzzysort from 'fuzzysort';
import moment, { Moment } from 'moment';
import { flow, getEnv } from 'mobx-state-tree';
import { Tag, ICategoryStore, ITag } from './category';
import { Api } from '../api';
import { Catalog, CatalogStore } from './catalog';
import { Kind, KindStore } from './kind';
import { assert } from './utils';

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
  description: types.optional(types.string, ''),
  minPipelinesVersion: types.optional(types.string, ''),
  rawURL: types.string,
  webURL: types.string,
  updatedAt: types.optional(updatedAt, '')
});

export const Resource = types
  .model('Resource', {
    id: types.number,
    name: types.identifier,
    catalog: types.reference(Catalog),
    kind: types.reference(Kind),
    latestVersion: types.reference(Version),
    displayVersion: types.reference(Version),
    tags: types.array(types.reference(Tag)), // ["1", "2"]
    rating: types.number,
    versions: types.array(types.reference(Version)),
    displayName: ''
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
    catalogs: types.optional(CatalogStore, {}),
    kinds: types.optional(KindStore, {}),
    sortBy: types.optional(types.enumeration(Object.values(SortByFields)), SortByFields.Unknown),
    tags: types.optional(types.map(Tag), {}),
    search: '',
    err: '',
    isLoading: true
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
    }
  }))

  .actions((self) => ({
    setLoading(l: boolean) {
      self.isLoading = l;
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
    }
  }))

  .actions((self) => ({
    versionInfo: flow(function* (resourceName: string) {
      try {
        self.setLoading(true);

        const { api } = self;
        const resource = self.resources.get(resourceName);
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
            if (self.resources.has(resourceName)) {
              const resource = self.resources.get(resourceName);
              assert(resource);
              resource.versions.push(v.id);
            }
          }
        });
      } catch (err) {
        self.err = err.toString();
      }
      self.setLoading(false);
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
          description: versionData.description,
          minPipelinesVersion: versionData.minPipelinesVersion,
          webURL: versionData.webURL,
          rawURL: versionData.rawURL,
          updatedAt: versionData.updatedAt
        };

        self.versions.put(version);
      } catch (err) {
        self.err = err.toString();
      }
      self.setLoading(false);
    }),

    load: flow(function* () {
      try {
        self.setLoading(true);

        const { api } = self;
        const json = yield api.resources();

        const kinds: string[] = json.data.map((r: IResource) => r.kind);
        kinds.forEach((k) => self.kinds.add(k));

        json.data.forEach((r: IResource) => {
          self.catalogs.add(r.catalog);
          self.versions.put(r.latestVersion);
        });

        // adding the tags to the store - normalized
        const tags: ITag[] = json.data.flatMap((item: IResource) => item.tags);

        tags.forEach((t) => (t != null ? self.tags.put(t) : null));

        const resources: IResource[] = json.data.map((r: IResource) => ({
          id: r.id,
          name: r.name,
          catalog: r.catalog.id,
          kind: r.kind,
          latestVersion: r.latestVersion.id,
          displayVersion: r.latestVersion.id,
          tags: r.tags != null ? r.tags.map((tag: ITag) => tag.id) : [],
          rating: r.rating,
          versions: [],
          displayName: r.latestVersion.displayName
        }));

        resources.forEach((r: IResource) => {
          r.versions.push(r.latestVersion);
          self.add(r);
        });
      } catch (err) {
        self.err = err.toString();
      }
      self.setLoading(false);
    })
  }))

  .actions((self) => ({
    afterCreate() {
      self.load();
    },
    setDisplayVersion(resourceName: string, versionId: string) {
      const resource = self.resources.get(resourceName);
      assert(resource);
      const version = self.versions.get(versionId);
      assert(version);
      if (version.id !== resource.displayVersion.id) {
        resource.displayVersion = version;
        self.versionUpdate(version.id);
      }
    }
  }))

  .views((self) => ({
    get filteredResources() {
      const { resources, kinds, catalogs, categories, search, sortBy } = self;
      const { selectedTags } = categories;

      let filteredItems: IResource[] = [];
      resources.forEach((r: IResource) => {
        const matchesKind = kinds.selected.size === 0 || kinds.selected.has(r.kind.name);
        const matchesCatalogs = catalogs.selected.size === 0 || catalogs.selected.has(r.catalog.id);
        const matchesTags = selectedTags.size === 0 || r.tags.some((t) => selectedTags.has(t.id));

        if (matchesKind && matchesCatalogs && matchesTags) {
          filteredItems.push(r);
        }
      });

      if (search.trim() !== '') {
        filteredItems = fuzzysort
          .go(search, filteredItems, { keys: ['name', 'displayName'] })
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
