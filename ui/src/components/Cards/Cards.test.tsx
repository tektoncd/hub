import React from 'react';
import { when } from 'mobx';
import { shallow } from 'enzyme';
import { FakeHub } from '../../api/testutil';
import { ResourceStore } from '../../store/resource';
import { CategoryStore } from '../../store/category';
import Cards from '.';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);

describe('Cards', () => {
  it('should render the resources on cards', (done) => {
    const store = ResourceStore.create({}, { api, categories: CategoryStore.create({}, { api }) });

    when(
      () => !store.isLoading,
      () => {
        const component = shallow(<Cards items={store.filteredResources} />);
        expect(component.debug()).toMatchSnapshot();

        done();
      }
    );
  });
});
