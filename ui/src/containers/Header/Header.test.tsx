import React from 'react';
import { mount } from 'enzyme';
import { when } from 'mobx';
import { BrowserRouter as Router } from 'react-router-dom';
import { Modal } from '@patternfly/react-core';
import Header from '.';
import Search from '../../containers/Search';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import { ActualDate, FakeDate } from '../../common/testutils';
import Icon from '../../components/Icon';
import { Icons } from '../../common/icons';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider, root } = createProviderAndStore(api);

// Assign a dummy date to global Date inorder to avoid test failue
FakeDate();

// Assign current Date once test is over
ActualDate();

describe('Header', () => {
  it('should render the header component and find Search component', () => {
    const component = mount(
      <Provider>
        <Router>
          <Header />
        </Router>
      </Provider>
    );

    expect(component.find(Search).length).toBe(1);
    expect(component.debug()).toMatchSnapshot();
  });

  it('should find Login', () => {
    const component = mount(
      <Provider>
        <Router>
          <Header />
        </Router>
      </Provider>
    );
    expect(component.find('span').slice(3).text()).toBe('Login');
  });

  it('should find Icons in header and their id', () => {
    const component = mount(
      <Provider>
        <Router>
          <Header />
        </Router>
      </Provider>
    );

    const icons = component.find(Icon);
    expect(icons.length).toBe(1);

    expect(icons.slice(0).props().id).toBe(Icons.Help);
  });

  it('should find the login Modal', () => {
    const component = mount(
      <Provider>
        <Router>
          <Header />
        </Router>
      </Provider>
    );

    expect(component.find(Modal).slice(3).props().className).toBe('hub-header-login__modal');
  });

  it('should authenticate', (done) => {
    const { user } = root;
    const code = {
      code: 'foo'
    };
    user.authenticate(code);
    when(
      () => {
        return !user.isLoading;
      },
      () => {
        setTimeout(() => {
          expect(user.isAuthenticated).toBe(true);

          done();
        }, 0);
      }
    );
  });

  it('should find Icon in header and it`s id', () => {
    const component = mount(
      <Provider>
        <Router>
          <Header />
        </Router>
      </Provider>
    );

    expect(component.find(Icon).length).toBe(1);
    expect(component.find(Icon).props().id).toBe(Icons.Help);
  });
});
