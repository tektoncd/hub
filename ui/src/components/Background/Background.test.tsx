import React from 'react';
import renderer from 'react-test-renderer';
import Background from '.';

it('should render the background component', () => {
  const component = renderer.create(<Background />).toJSON();
  expect(component).toMatchSnapshot();
});
