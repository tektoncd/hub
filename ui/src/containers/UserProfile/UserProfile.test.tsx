import React from 'react';
import { mount } from 'enzyme';
import { Avatar, Dropdown } from '@patternfly/react-core';
import UserProfile from '.';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import { assert } from '../../store/utils';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider } = createProviderAndStore(api);

describe('UserProfile', () => {
  it('renders component correctly', () => {
    const component = mount(
      <Provider>
        <UserProfile />
      </Provider>
    );
    expect(component.find(Avatar).length).toBe(1);

    const dropdownItems = component.find(Dropdown).props().dropdownItems;
    assert(dropdownItems);
    expect(dropdownItems.length).toBe(2);
  });
});
