import { Card } from '@patternfly/react-core';
import { mount, shallow } from 'enzyme';
import React from 'react';
import renderer from 'react-test-renderer';
import GitHubLogin from 'react-github-login';
import Authentication from '.';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider } = createProviderAndStore(api);

describe('Authentication', () => {
  it('renders component correctly', () => {
    const tree = renderer
      .create(
        <Provider>
          <Authentication />
        </Provider>
      )
      .toJSON();
    expect(tree).toMatchSnapshot();
  });

  it('it can find card and GitHubLogin', () => {
    const component = shallow(<Authentication />);
    expect(component.find(Card).length).toBe(1);
    expect(component.find(GitHubLogin).length).toBe(1);
  });

  it('it can test github login', () => {
    window.open = jest.fn();
    const windowOpenSpy = jest.spyOn(window, 'open');
    const clientId = 'foo';
    const redirectUri = 'http://foo.test/auth/github';

    const component = <GitHubLogin clientId={clientId} redirectUri={redirectUri} />;
    const wrapper = mount(component);

    wrapper.find('button').simulate('click');

    const query = `client_id=${clientId}&scope=user:email&redirect_uri=${redirectUri}`;

    expect(wrapper.instance().popup.url).toBe(`https://github.com/login/oauth/authorize?${query}`);
    expect(windowOpenSpy).toHaveBeenCalled();
  });
});
