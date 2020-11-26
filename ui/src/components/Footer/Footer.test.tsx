import React from 'react';
import renderer from 'react-test-renderer';
import Footer from '.';

it('should render the footer component', () => {
  const component = renderer.create(<Footer />).toJSON();
  expect(component).toMatchSnapshot();
});
