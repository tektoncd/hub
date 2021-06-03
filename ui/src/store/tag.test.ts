import { getSnapshot } from 'mobx-state-tree';
import { Tag, TagStore } from './tag';
import { assert } from './utils';

describe('Store Object', () => {
  it('can create a tag object', () => {
    const tag = Tag.create({
      id: 1,
      name: 'aws'
    });

    expect(tag.name).toBe('aws');
  });

  it('creates a tag store', (done) => {
    const store = TagStore.create({});

    const item = Tag.create({
      id: 10,
      name: 'cli'
    });

    store.add(item);

    expect(store.values[0].name).toBe('cli');

    const tags = store.items.get('cli');
    assert(tags);

    expect(tags.name).toBe('cli');
    expect(getSnapshot(store.items)).toMatchSnapshot();

    done();
  });
});
