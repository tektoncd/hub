import React from 'react';
import renderer from 'react-test-renderer';
import Rating from '.';

it('should render Ratings component', () => {
  const component = renderer.create(<Rating />).toJSON();
  expect(component).toMatchSnapshot();
});
