import { types, Instance } from 'mobx-state-tree';
import fuzzysort from 'fuzzysort';
import moment, { Moment } from 'moment';
import { flow, getEnv } from 'mobx-state-tree';
import { Tag, ICategoryStore, ITag } from './category';
import { Api } from '../api';
import { Catalog, CatalogStore } from './catalog';
import { Kind, KindStore } from './kind';

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
  displayName: types.string,
  description: types.string,
  minPipelinesVersion: types.string,
  rawURL: types.string,
  webURL: types.string,
  updatedAt: updatedAt
});

export const Resource = types.model('Resource', {
  id: types.identifierNumber,
  name: types.optional(types.string, ''),
  catalog: types.reference(Catalog),
  kind: types.reference(Kind),
  latestVersion: types.reference(Version),
  tags: types.array(types.reference(Tag)), // ["1", "2"]
  rating: types.number,
  versions: types.array(types.reference(Version)),
  displayName: ''
});

export type IResource = Instance<typeof Resource>;
export type IVersion = Instance<typeof Version>;

export enum sortByFields {
  Name = 'Name',
  Rating = 'Rating'
}

export const ResourceStore = types
  .model('ResourceStore', {
    resources: types.map(Resource),
    versions: types.map(Version),
    catalogs: types.optional(CatalogStore, {}),
    kinds: types.optional(KindStore, {}),
    sortBy: types.optional(types.enumeration(Object.values(sortByFields)), sortByFields.Name),
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
      const key: sortByFields = sortByFields[field as keyof typeof sortByFields];
      self.sortBy = key;
    }
  }))

  .actions((self) => ({
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

        const resources: IResource[] = json.data.map((r: IResource) => ({
          id: r.id,
          name: r.name,
          catalog: r.catalog.id,
          kind: r.kind,
          latestVersion: r.latestVersion.id,
          tags: r.tags.map((tag: ITag) => tag.id),
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
    }
  }))

  .views((self) => ({
    get filteredResources() {
      const { resources, kinds, catalogs, categories, search } = self;
      const { selectedTags } = categories;

      const filtered: IResource[] = [];
      resources.forEach((r: IResource) => {
        const matchesKind = kinds.selected.size === 0 || kinds.selected.has(r.kind.name);
        const matchesCatalogs = catalogs.selected.size === 0 || catalogs.selected.has(r.catalog.id);
        const matchesTags = selectedTags.size === 0 || r.tags.some((t) => selectedTags.has(t.id));

        if (matchesKind && matchesCatalogs && matchesTags) {
          filtered.push(r);
        }
      });

      if (search.trim() !== '') {
        return fuzzysort
          .go(search, filtered, { keys: ['name', 'displayName'] })
          .map((resource: Fuzzysort.KeysResult<IResource>) => resource.obj);
      }
      return filtered;
    }
  }));

export type IResourceStore = Instance<typeof ResourceStore>;
