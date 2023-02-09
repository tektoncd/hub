import React from 'react';
import { mount } from 'enzyme';
import { when } from 'mobx';
import { EmptyState, GalleryItem, Button } from '@patternfly/react-core';
import { BrowserRouter as Router } from 'react-router-dom';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import Cards from '../../components/Cards';
import Resources from '.';
import { SortByFields } from '../../store/resource';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider, root } = createProviderAndStore(api);

beforeEach(() => {
  global.Date.now = jest.fn(() => new Date('2020-12-22T10:20:30Z').getTime());
});

afterEach(() => {
  global.Date = Date;
});

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
    const { resources } = root;
    when(
      () => {
        return !resources.isLoading;
      },
      () => {
        when(
          () => {
            return !resources.isLoading;
          },
          () => {
            setTimeout(() => {
              const component = mount(
                <Provider>
                  <Router>
                    <Resources />
                  </Router>
                </Provider>
              );
              component.update();

              expect(component.find(EmptyState).length).toBe(1);
            }, 1000);
          }
        );
        done();
      }
    );
  });

  it('should find Clear All Filters button in the EmptyState', (done) => {
    const { resources } = root;
    when(
      () => {
        return !resources.isLoading;
      },
      () => {
        when(
          () => {
            return !resources.isLoading;
          },
          () => {
            setTimeout(() => {
              const component = mount(
                <Provider>
                  <Router>
                    <Resources />
                  </Router>
                </Provider>
              );
              component.update();

              component.update();
              const r = component.find(EmptyState);

              expect(r.find(Button).length).toEqual(1);
            }, 1000);
          }
        );
        done();
      }
    );
  });
});
