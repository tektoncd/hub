import React from 'react';
import { mount } from 'enzyme';
import Sort from '.';
import { Dropdown } from '@patternfly/react-core';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import { assert } from '../../store/utils';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider } = createProviderAndStore(api);

describe('Sort DropDown', () => {
  it('should render the dropdown', () => {
    const component = mount(
      <Provider>
        <Sort />
      </Provider>
    );

    const item = component.find(Dropdown).props().dropdownItems;
    assert(item);

    expect(item.length).toBe(2);
    expect(item[0].key).toBe('Name');
    expect(component.debug()).toMatchSnapshot();
  });
});
