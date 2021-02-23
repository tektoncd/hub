import React from 'react';
import { shallow } from 'enzyme';
import { PageNotFound } from './index';

jest.mock('react-router-dom', () => {
  return {
    useHistory: () => {
      return {
        history: ''
      };
    }
  };
});

describe('PageNotFound Component', () => {
  it('should render PageNotFound component', () => {
    const component = shallow(<PageNotFound />);
    expect(component.debug()).toMatchSnapshot();
  });
});
