import React from 'react';
import { mount } from 'enzyme';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore, createProvider } from '../../store/root';
import LeftPane from '../../components/LeftPane';
import App from '.';
import { FakeDate, ActualDate } from '../../common/testutils';
import Details from '../Details';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider } = createProviderAndStore(api);

// Assign a dummy date to global Date inorder to avoid snapshaot updation
FakeDate();

// Assign current Date once test is over
ActualDate();

describe('App', () => {
  it('should render the component correctly and match the snapshot', (done) => {
    const app = mount(
      <Provider>
        <App />
      </Provider>
    );

    expect(app.debug()).toMatchSnapshot();
    done();
  });

  it('should find the component and match the count', () => {
    const Provider = createProvider();
    const component = mount(
      <Provider>
        <App />
      </Provider>
    );
    expect(component.find(LeftPane).length).toEqual(1);
  });
});
