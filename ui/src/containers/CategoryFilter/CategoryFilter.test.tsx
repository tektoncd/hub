import React from 'react';
import { mount } from 'enzyme';
import CategoryFilter from '.';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import { when } from 'mobx';
import { Checkbox } from '@patternfly/react-core';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider, root } = createProviderAndStore(api);

describe('CategoryFilter', () => {
  it('finds the filter component and matches the count', (done) => {
    const component = mount(
      <Provider>
        <CategoryFilter />
      </Provider>
    );

    const { categories } = root;
    when(
      () => {
        return !categories.isLoading;
      },
      () => {
        setTimeout(() => {
          expect(categories.count).toBe(5);
          component.update();

          const cf = component.find(CategoryFilter);
          expect(cf.length).toEqual(1);

          expect(component.debug()).toMatchSnapshot();

          expect(component.find(Checkbox).length).toEqual(5);

          done();
        }, 0);
      }
    );
  });
});
