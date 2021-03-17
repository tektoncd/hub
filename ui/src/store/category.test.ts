import { when } from 'mobx';
import { getSnapshot } from 'mobx-state-tree';
import { FakeHub } from '../api/testutil';
import { CategoryStore, Category, Tag } from './category';
import { assert } from './utils';

const TESTDATA_DIR = `${__dirname}/testdata`;
const api = new FakeHub(TESTDATA_DIR);

describe('Store Object', () => {
  it('can create a tag object', () => {
    const store = Tag.create({
      id: 1,
      name: 'cli'
    });
    expect(store.name).toBe('cli');
  });
  it('can create a category object', () => {
    const category = Category.create({
      id: 1,
      name: 'test',
      tags: ['1']
    });

    expect(category.name).toBe('test');
    expect(category.tags.length).toBe(1);
  });
});

describe('Store functions', () => {
  it('can create a category store', (done) => {
    const store = CategoryStore.create({}, { api });
    expect(store.count).toBe(0);
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        expect(store.count).toBe(5);
        expect(store.isLoading).toBe(false);

        expect(getSnapshot(store)).toMatchSnapshot();

        done();
      }
    );
  });

  it('can toggle the selected category', (done) => {
    const store = CategoryStore.create({}, { api });
    expect(store.count).toBe(0);
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.count).toBe(5);
        expect(store.isLoading).toBe(false);

        const categories = store.items.get('1');
        assert(categories);
        categories.toggle();

        expect(categories.selected).toBe(true);

        done();
      }
    );
  });

  it('can clear all the categories', (done) => {
    const store = CategoryStore.create({}, { api });
    expect(store.count).toBe(0);

    when(
      () => !store.isLoading,
      () => {
        expect(store.count).toBe(5);
        expect(store.isLoading).toBe(false);

        // Gets the category with id as 1
        const c1 = store.items.get('1');
        assert(c1);
        c1.toggle();

        // Gets the category with id as 2
        const c2 = store.items.get('2');
        assert(c2);
        c2.toggle();

        store.clearSelected();
        expect(c1.selected).toBe(false);

        done();
      }
    );
  });

  it('can return the tags for the categories which are selected', (done) => {
    const store = CategoryStore.create({}, { api });
    expect(store.count).toBe(0);
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.count).toBe(5);
        expect(store.isLoading).toBe(false);

        // Gets the category with id as 1
        const c1 = store.items.get('1');
        assert(c1);
        c1.toggle();

        // Gets the category with id as 2
        const c2 = store.items.get('2');
        assert(c2);
        c2.toggle();

        expect(store.selectedTags.size).toBe(2);

        done();
      }
    );
  });

  it('can toggle the category by name', (done) => {
    const store = CategoryStore.create({}, { api });
    expect(store.count).toBe(0);
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.count).toBe(5);
        expect(store.isLoading).toBe(false);

        store.toggleByName('Build Tools');

        const categories = store.items.get('1');
        assert(categories);
        expect(categories.selected).toBe(true);

        done();
      }
    );
  });

  it('can return the all selected catgories in a list', (done) => {
    const store = CategoryStore.create({}, { api });
    expect(store.count).toBe(0);
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.count).toBe(5);
        expect(store.isLoading).toBe(false);

        // Gets the category with id as 1
        const c1 = store.items.get('1');
        assert(c1);
        c1.toggle();

        // Gets the category with id as 2
        const c2 = store.items.get('2');
        assert(c2);
        c2.toggle();

        expect(store.selectedByName.sort()).toEqual(['Build Tools', 'CLI'].sort());

        done();
      }
    );
  });
});
