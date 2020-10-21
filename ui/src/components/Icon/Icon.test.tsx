import React from 'react';
import { shallow } from 'enzyme';
import Icon from '.';
import { Icons } from '../../common/icons';

describe('Icon Component', () => {
  it('should render icon for Task Kind', () => {
    const component = shallow(<Icon id={Icons.Build} size="sm" label="Task" />);
    expect(component.debug()).toMatchSnapshot();
    expect(component.find('BuildIcon[label="Task"]').length).toEqual(1);
  });
  it('should render icon for Pipeline Kind', () => {
    const component = shallow(<Icon id={Icons.Domain} size="md" label="Pipeline" />);
    expect(component.debug()).toMatchSnapshot();
    expect(component.find('DomainIcon[label="Pipeline"]').length).toEqual(1);
  });
  it('should render icon for Official Catalog', () => {
    const component = shallow(<Icon id={Icons.Cat} size="lg" label="Official" />);
    expect(component.debug()).toMatchSnapshot();
    expect(component.find('CatIcon[label="Official"]').length).toEqual(1);
  });
  it('should render icon for Verified Catalog', () => {
    const component = shallow(<Icon id={Icons.Certificate} size="xl" label="Verified" />);
    expect(component.debug()).toMatchSnapshot();
    expect(component.find('CertificateIcon[label="Verified"]').length).toEqual(1);
  });
  it('should render icon for Community Catalog', () => {
    const component = shallow(<Icon id={Icons.User} size="sm" label="Community" />);
    expect(component.debug()).toMatchSnapshot();
    expect(component.find('UserIcon[label="Community"]').length).toEqual(1);
  });
  it('should render icon for Category Filter', () => {
    const component = shallow(<Icon id={Icons.Unknown} size="sm" label="CLI" />);
    expect(component.debug()).toMatchSnapshot();
  });
});
