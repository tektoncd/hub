import React from 'react';
import renderer from 'react-test-renderer';
import Search from '.';

it('should render the search component', () => {
  const tree = renderer.create(<Search />).toJSON();
  expect(tree).toMatchSnapshot();
});
