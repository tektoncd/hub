import React from 'react';
import { mount } from 'enzyme';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import { when } from 'mobx';
import { Checkbox } from '@patternfly/react-core';
import CatalogFilter from '.';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider, root } = createProviderAndStore(api);

describe('CatalogFilter', () => {
  it('finds the filter component and matches the count', (done) => {
    const component = mount(
      <Provider>
        <CatalogFilter />
      </Provider>
    );

    const { resources } = root;
    when(
      () => {
        return !resources.isLoading;
      },
      () => {
        setTimeout(() => {
          const catalogs = resources.catalogs;
          expect(catalogs.items.size).toBe(2);
          component.update();

          const cf = component.find(CatalogFilter);
          expect(cf.length).toEqual(1);

          expect(component.debug()).toMatchSnapshot();

          expect(component.find(Checkbox).length).toEqual(2);

          done();
        }, 0);
      }
    );
  });
});
