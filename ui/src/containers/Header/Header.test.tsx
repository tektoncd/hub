import React from 'react';
import { mount } from 'enzyme';
import { when } from 'mobx';
import { BrowserRouter as Router } from 'react-router-dom';
import Header from '.';
import Search from '../../containers/Search';
import UserProfile from '../UserProfile';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider, root } = createProviderAndStore(api);

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
    expect(component.find('span').text()).toBe('Login');
  });

  it('should render user profile', (done) => {
    const component = mount(
      <Provider>
        <Router>
          <Header />
        </Router>
      </Provider>
    );
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
          component.update();
          expect(component.find(UserProfile).length).toEqual(1);
          done();
        }, 0);
      }
    );
  });
});
