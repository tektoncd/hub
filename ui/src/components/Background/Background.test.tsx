import React from 'react';
import { shallow } from 'enzyme';
import Background from '.';

it('should render the background component', () => {
  const component = shallow(<Background />);
  expect(component.debug()).toMatchSnapshot();
});
