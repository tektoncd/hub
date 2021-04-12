import React from 'react';
import { FakeHub } from '../../api/testutil';
import { when } from 'mobx';
import { shallow } from 'enzyme';
import Filter from '.';
import { CategoryStore } from '../../store/category';
import { ResourceStore } from '../../store/resource';
import { CatalogStore } from '../../store/catalog';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);

describe('Filter component', () => {
  it('should render the component', (done) => {
    const store = CategoryStore.create({}, { api });

    when(
      () => !store.isLoading,
      () => {
        expect(store.count).toBe(5);
        expect(store.isLoading).toBe(false);

        const component = shallow(<Filter store={store} header="Categories" />);
        expect(component).toMatchSnapshot();

        done();
      }
    );
  });

  it('checks the checkbox render count for categories', (done) => {
    const store = CategoryStore.create({}, { api });

    when(
      () => !store.isLoading,
      () => {
        const component = shallow(<Filter store={store} header="Categories" />);
        expect(component.find('Checkbox[aria-label="controlled checkbox"]').length).toEqual(5);

        done();
      }
    );
  });

  it('finds the button and matches the count', (done) => {
    const store = CategoryStore.create({}, { api });

    when(
      () => !store.isLoading,
      () => {
        const component = shallow(<Filter store={store} header="Categories" />);

        expect(component.find('Button[variant="plain"]').length).toEqual(1);

        done();
      }
    );
  });

  it('checks the checkbox render count for catalogs', (done) => {
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
        const component = shallow(<Filter store={store.catalogs} header="Catalogs" />);
        expect(component.find('Checkbox[aria-label="controlled checkbox"]').length).toEqual(2);

        done();
      }
    );
  });

  it('checks the checkbox render count for kinds', (done) => {
    const store = ResourceStore.create({}, { api, categories: CategoryStore.create({}, { api }) });

    when(
      () => !store.isLoading,
      () => {
        const component = shallow(<Filter store={store.kinds} header="Kind" />);
        expect(component.find('Checkbox[aria-label="controlled checkbox"]').length).toEqual(2);

        done();
      }
    );
  });
});
