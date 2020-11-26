import React from 'react';
import renderer from 'react-test-renderer';
import Details from '.';

it('should render the details component', () => {
  const component = renderer.create(<Details />).toJSON();
  expect(component).toMatchSnapshot();
});
