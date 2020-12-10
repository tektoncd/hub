import React from 'react';
import { mount } from 'enzyme';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import Header from '.';
import Search from '../../containers/Search';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider } = createProviderAndStore(api);

it('should render the header component and finds Search component', () => {
  const component = mount(
    <Provider>
      <Header />
    </Provider>
  );

  expect(component.find(Search).length).toBe(1);
  expect(component.debug()).toMatchSnapshot();
});
