import React from 'react';
import { shallow } from 'enzyme';
import App from '.';

test('renders learn react link', () => {
  const app = shallow(<App />);
  expect(app.debug()).toMatchSnapshot();
});
