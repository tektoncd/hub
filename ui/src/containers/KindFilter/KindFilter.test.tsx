import React from 'react';
import { mount } from 'enzyme';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import { when } from 'mobx';
import { Checkbox } from '@patternfly/react-core';
import KindFilter from '.';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider, root } = createProviderAndStore(api);

describe('CategoryFilter', () => {
  it('finds the filter component and matches the count', (done) => {
    const component = mount(
      <Provider>
        <KindFilter />
      </Provider>
    );

    const { resources } = root;
    when(
      () => {
        return !resources.isLoading;
      },
      () => {
        setTimeout(() => {
          const kinds = resources.kinds;
          expect(kinds.items.size).toBe(2);
          component.update();

          const kf = component.find(KindFilter);
          expect(kf.length).toEqual(1);

          expect(component.debug()).toMatchSnapshot();

          expect(component.find(Checkbox).length).toEqual(2);

          done();
        }, 0);
      }
    );
  });
});
