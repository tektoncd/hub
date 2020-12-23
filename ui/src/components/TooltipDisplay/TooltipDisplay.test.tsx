import React from 'react';
import { shallow } from 'enzyme';
import { Tooltip } from '@patternfly/react-core';
import { Icons } from '../../common/icons';
import TooltipDisplay from '.';
import { assert } from '../../store/utils';

describe('TooltipDisplay component', () => {
  it('should render Tooltip for resource kind icon', () => {
    const component = shallow(<TooltipDisplay id={Icons.Build} name="Task" />);
    expect(component.debug()).toMatchSnapshot();

    const propsChildren = component.find(Tooltip).props().children;
    assert(propsChildren);
    expect(propsChildren.props['label']).toBe('Task');
  });

  it('should render Tooltip for resource catalog icon', () => {
    const component = shallow(<TooltipDisplay id={Icons.User} name="Tekton" />);
    expect(component.debug()).toMatchSnapshot();

    const propsChildren = component.find(Tooltip).props().children;
    assert(propsChildren);
    expect(propsChildren.props['label']).toBe('Tekton');
  });

  it('should render Tooltip for resource Avergae Rating icon', () => {
    const component = shallow(<TooltipDisplay id={Icons.Star} name="Avergae Rating" />);
    expect(component.debug()).toMatchSnapshot();

    const propsChildren = component.find(Tooltip).props().children;
    assert(propsChildren);
    expect(propsChildren.props['label']).toBe('Avergae Rating');
  });
});
