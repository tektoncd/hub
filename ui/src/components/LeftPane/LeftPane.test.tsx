import React from 'react';
import LeftPane from '.';
import { shallow } from 'enzyme';

import KindFilter from '../../containers/KindFilter';
import CatalogFilter from '../../containers/CatalogFilter';
import CategoryFilter from '../../containers/CategoryFilter';
import PlatformFilter from '../../containers/PlatformFilter';

describe('LeftPane', () => {
  it('should find the components and match the count', () => {
    const component = shallow(<LeftPane />);

    expect(component.debug()).toMatchSnapshot();

    expect(component.find(KindFilter).length).toEqual(1);
    expect(component.find(CatalogFilter).length).toEqual(1);
    expect(component.find(PlatformFilter).length).toEqual(1);
    expect(component.find(CategoryFilter).length).toEqual(1);
  });
});
