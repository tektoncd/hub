import React from 'react';
import { AlertVariant } from '@patternfly/react-core';
import { shallow } from 'enzyme';
import { IError } from '../../store/auth';
import AlertDisplay from '.';

it('should render the Alert component', () => {
  const error: IError = {
    status: 404,
    customMessage: 'Test message',
    serverMessage: 'Test server message'
  };
  const component = shallow(<AlertDisplay message={error} alertVariant={AlertVariant.danger} />);
  expect(component.debug()).toMatchSnapshot();
});
