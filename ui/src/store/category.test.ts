import { when } from 'mobx';
import { getSnapshot } from 'mobx-state-tree';
import { FakeHub } from '../api/testutil';
import { CategoryStore, Category } from './category';

const TESTDATA_DIR = `${__dirname}/testdata`;
const api = new FakeHub(TESTDATA_DIR);

describe('Category Object', () => {
  it('can create a category object', () => {
    const category = Category.create({
      id: 1,
      name: 'test',
      tags: [
        {
          id: 1,
          name: 'test-category'
        }
      ]
    });

    expect(category.name).toBe('test');
    expect(category.id).toBe(1);
    expect(category.tags[0].name).toBe('test-category');
  });
});

describe('Store functions', () => {
  it('can create a store', (done) => {
    const store = CategoryStore.create({}, { api });
    expect(store.count).toBe(0);
    expect(store.isLoading).toBe(true);
    when(
      () => !store.isLoading,
      () => {
        expect(store.count).toBe(5);
        expect(store.isLoading).toBe(false);

        expect(store.list[0].id).toBe(1);
        expect(store.list[0].name).toBe('Build Tools');
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

        store.list[0].toggle();

        expect(store.list[0].selected).toBe(true);

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

        store.list[0].toggle();
        store.list[1].toggle();

        const tags = store.tags;
        expect(tags[0]).toBe('build-tool');

        done();
      }
    );
  });

  it('clears all the selected categories', (done) => {
    const store = CategoryStore.create({}, { api });
    expect(store.count).toBe(0);
    expect(store.isLoading).toBe(true);

    when(
      () => !store.isLoading,
      () => {
        expect(store.count).toBe(5);
        expect(store.isLoading).toBe(false);

        store.list[0].toggle();
        store.list[2].toggle();
        store.clear();

        expect(store.list).toEqual(
          expect.arrayContaining([
            expect.objectContaining({
              selected: false
            })
          ])
        );

        done();
      }
    );
  });
});
