import { ReactComponent } from '*.svg';
import { FormSelectOption, TextInput } from '@patternfly/react-core';
import { mount, shallow } from 'enzyme';
import { when } from 'mobx';
import React from 'react';
import { act } from 'react-dom/test-utils';
import Search from '.';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider, root } = createProviderAndStore(api);

jest.mock('react-router-dom', () => ({
  useNavigate: () => ({
    replace: jest.fn()
  })
}));

describe('search resources', () => {
  it('should render the search component', () => {
    const component = shallow(
      <Provider>
        <Search />
      </Provider>
    );
    expect(component.debug()).toMatchSnapshot();
  });

  it('Should set value to store when input is changed', (done) => {
    const component = mount(
      <Provider>
        <Search />
      </Provider>
    );

    const { resources } = root;
    when(
      () => {
        return !resources.isLoading;
      },
      () => {
        act(() => {
          component.find(TextInput).simulate('change');
        });
        const resourceName = jest.fn(() => 'golang');

        expect(resourceName()).toBe('golang');
        expect(resourceName).toHaveBeenCalledTimes(1);

        done();
      }
    );
  });
});
