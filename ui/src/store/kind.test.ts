import { Kind, KindStore } from './kind';
import { getSnapshot } from 'mobx-state-tree';
import { assert } from './utils';
import { Icons } from '../common/icons';

describe('Store Object', () => {
  it('can create a kind object', () => {
    const store = Kind.create({
      name: 'Task'
    });

    expect(store.name).toBe('Task');
  });

  it('creates a kind store', (done) => {
    const store = KindStore.create({});

    const item = Kind.create({
      name: 'Task'
    });

    store.add(item.name);

    const kinds = store.items.get('Task');
    assert(kinds);

    expect(kinds.name).toBe('Task');
    expect(getSnapshot(store.items)).toMatchSnapshot();

    done();
  });

  it('should toggle a selected kind', (done) => {
    const store = KindStore.create({});

    const item = Kind.create({
      name: 'Task'
    });

    store.add(item.name);

    const kinds = store.items.get('Task');
    assert(kinds);
    kinds.toggle();

    expect(store.selected.size).toBe(1);
    expect(kinds.selected).toBe(true);

    done();
  });

  it('should clear all the selected kind', (done) => {
    const store = KindStore.create({});

    const item = {
      name: 'Task'
    };

    store.add(item.name);

    const kinds = store.items.get('Task');
    assert(kinds);
    kinds.toggle();

    store.clearSelected();
    expect(kinds.selected).toBe(false);

    done();
  });

  it('should get an icon for kind', (done) => {
    const store = KindStore.create({});

    const item = {
      name: 'Task'
    };

    store.add(item.name);

    const kind = store.items.get('Task');
    assert(kind);
    expect(kind.icon).toBe(Icons.Build);

    done();
  });
});
