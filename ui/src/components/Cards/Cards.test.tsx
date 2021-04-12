import React from 'react';
import { when } from 'mobx';
import { shallow } from 'enzyme';
import { FakeHub } from '../../api/testutil';
import { ResourceStore } from '../../store/resource';
import { CategoryStore } from '../../store/category';
import { CatalogStore } from '../../store/catalog';

import Cards from '.';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);

beforeEach(() => {
  global.Date.now = jest.fn(() => new Date('2020-12-22T10:20:30Z').getTime());
});

afterEach(() => {
  global.Date = Date;
});

describe('Cards', () => {
  it('should render the resources on cards', (done) => {
    const store = ResourceStore.create(
      {},
      {
        api,
        catalogs: CatalogStore.create({}, { api }),
        categories: CategoryStore.create({}, { api })
      }
    );

    when(
      () => !store.isLoading,
      () => {
        const component = shallow(<Cards items={store.filteredResources} />);
        expect(component.debug()).toMatchSnapshot();

        done();
      }
    );
  });
});
