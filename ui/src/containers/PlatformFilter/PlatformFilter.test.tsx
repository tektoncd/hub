import React from 'react';
import { mount } from 'enzyme';
import { when } from 'mobx';
import { Checkbox } from '@patternfly/react-core';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import PlatformFilter from '.';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider, root } = createProviderAndStore(api);

describe('PlatformFilter', () => {
  it('finds the filter component and matches the count', (done) => {
    const component = mount(
      <Provider>
        <PlatformFilter />
      </Provider>
    );

    const { resources } = root;
    when(
      () => {
        return !resources.isLoading;
      },
      () => {
        setTimeout(() => {
          const platforms = resources.platforms;
          expect(platforms.items.size).toBe(4);
          component.update();

          const pf = component.find(PlatformFilter);
          expect(pf.length).toEqual(1);

          expect(component.debug()).toMatchSnapshot();

          expect(component.find(Checkbox).length).toEqual(4);

          done();
        }, 0);
      }
    );
  });
});
