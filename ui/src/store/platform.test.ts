import { Platform, PlatformStore } from './platform';
import { getSnapshot } from 'mobx-state-tree';
import { assert } from './utils';

describe('Store Object', () => {
  it('can create a platform object', () => {
    const platform = Platform.create({
      id: 1,
      name: 'linux/amd64'
    });

    expect(platform.name).toBe('linux/amd64');
  });

  it('creates a platform store', (done) => {
    const store = PlatformStore.create({});

    const item = Platform.create({
      id: 1,
      name: 'linux/amd64'
    });

    store.add(item);

    const platform = store.items.get('1');
    assert(platform);

    expect(platform.name).toBe('linux/amd64');
    expect(getSnapshot(store.items)).toMatchSnapshot();

    done();
  });

  it('should toggle a selected platform', (done) => {
    const store = PlatformStore.create({});

    const item = Platform.create({
      id: 1,
      name: 'linux/amd64'
    });

    store.add(item);

    const platform = store.items.get('1');
    assert(platform);
    platform.toggle();

    expect(store.selected.size).toBe(1);
    expect(platform.selected).toBe(true);

    done();
  });

  it('should clear all the selected platforms', (done) => {
    const store = PlatformStore.create({});

    const item = Platform.create({
      id: 1,
      name: 'linux/amd64',
      selected: false
    });

    store.add(item);

    const platform = store.items.get('1');
    assert(platform);
    platform.toggle();

    store.clearSelected();
    expect(platform.selected).toBe(false);

    done();
  });
});
