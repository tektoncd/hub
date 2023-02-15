import React from 'react';
import { shallow } from 'enzyme';
import HomePage from '.';

it('should render the home page component', () => {
  const component = shallow(<HomePage />);
  expect(component.debug()).toMatchSnapshot();
});
