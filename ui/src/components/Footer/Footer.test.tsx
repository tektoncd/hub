import React from 'react';
import { shallow } from 'enzyme';
import Footer from '.';
import { FakeDate, ActualDate } from '../../common/testutils';

// Assign a dummy date to global Date inorder to avoid snapshaot updation
FakeDate();

// Assign current Date once test is over
ActualDate();

it('should render the footer component', () => {
  const component = shallow(<Footer />);
  expect(component.debug()).toMatchSnapshot();
});
