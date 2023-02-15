import React from 'react';
import { shallow } from 'enzyme';
import ResourceDetailsPage from '.';

it('should render the resource details page  component', () => {
  const component = shallow(<ResourceDetailsPage />);
  expect(component.debug()).toMatchSnapshot();
});
