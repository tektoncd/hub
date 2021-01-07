import React from 'react';
import { mount } from 'enzyme';
import { when } from 'mobx';
import { EmptyState, GalleryItem } from '@patternfly/react-core';
import { BrowserRouter as Router } from 'react-router-dom';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import Cards from '../../components/Cards';
import Resources from '.';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider, root } = createProviderAndStore(api);

describe('Resource Component', () => {
  it('should render the resources component', (done) => {
    const component = mount(
      <Provider>
        <Router>
          <Resources />
        </Router>
      </Provider>
    );

    const { resources } = root;
    when(
      () => {
        return !resources.isLoading;
      },
      () => {
        setTimeout(() => {
          const resource = resources.filteredResources;
          expect(resource.length).toBe(7);

          component.update();

          const r = component.find(Resources);
          expect(r.length).toEqual(1);

          expect(component.debug()).toMatchSnapshot();

          const c = component.find(Cards);
          expect(c.find(GalleryItem).length).toBe(7);

          done();
        }, 0);
      }
    );
  });

  it('should find EmptyState if filtered does not match to any resources', (done) => {
    const component = mount(
      <Provider>
        <Router>
          <Resources />
        </Router>
      </Provider>
    );

    const { resources } = root;
    when(
      () => {
        return !resources.isLoading;
      },
      () => {
        setTimeout(() => {
          resources.setSearch('asdf');
          const resource = resources.filteredResources;
          expect(resource.length).toBe(0);

          component.update();
          const r = component.find(EmptyState);
          expect(r.length).toEqual(1);

          done();
        }, 0);
      }
    );
  });
});
