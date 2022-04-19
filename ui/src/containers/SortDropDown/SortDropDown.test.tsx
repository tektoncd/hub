import React from 'react';
import { mount } from 'enzyme';
import { Select } from '@patternfly/react-core';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import { assert } from '../../store/utils';
import Sort from '.';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider } = createProviderAndStore(api);

describe('Sort DropDown', () => {
  it('it should render the selected dropdrown', () => {
    const component = mount(
      <Provider>
        <Sort />
      </Provider>
    );

    const item = component.find(Select);

    expect(item.length).toBe(1);
    expect(item.props().placeholderText).toBe('Sort By');
    expect(item.props().selections).toBe('');

    const keys = item.props().children;
    assert(keys);
    expect(keys.length).toBe(3);
  });
});
