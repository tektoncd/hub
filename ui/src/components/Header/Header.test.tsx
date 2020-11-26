import React from 'react';
import renderer from 'react-test-renderer';
import Header from '.';

it('should render the header component', () => {
  const component = renderer.create(<Header />).toJSON();
  expect(component).toMatchSnapshot();
});
