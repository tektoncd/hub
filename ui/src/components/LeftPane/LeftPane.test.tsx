import React from 'react';
import LeftPane from '.';
import { mount } from 'enzyme';

import KindFilter from '../../containers/KindFilter';
import CatalogFilter from '../../containers/CatalogFilter';
import CategoryFilter from '../../containers/CategoryFilter';
import PlatformFilter from '../../containers/PlatformFilter';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider } = createProviderAndStore(api);

describe('LeftPane', () => {
  it('should find the components and match the count', () => {
    const component = mount(
      <Provider>
        <LeftPane />
      </Provider>
    );

    expect(component.debug()).toMatchSnapshot();

    expect(component.find(KindFilter).length).toEqual(1);
    expect(component.find(CatalogFilter).length).toEqual(1);
    expect(component.find(PlatformFilter).length).toEqual(1);
    expect(component.find(CategoryFilter).length).toEqual(1);
  });
});
