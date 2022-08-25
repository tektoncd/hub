import React from 'react';
import renderer from 'react-test-renderer';
import Background from '.';

fit('should render the background component', () => {
  const component = renderer.create(<Background />).toJSON();
  expect(component).toMatchSnapshot();
});
