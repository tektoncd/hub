import React from 'react';
import { Card, CardActions, Dropdown } from '@patternfly/react-core';
import { mount } from 'enzyme';
import { when } from 'mobx';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import BasicDetails from '.';
import { assert } from '../../store/utils';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider, root } = createProviderAndStore(api);

jest.mock('react-router-dom', () => {
  return {
    useParams: () => {
      return {
        name: 'buildah'
      };
    }
  };
});

it('should render the BasicDetails component', (done) => {
  const { resources } = root;
  when(
    () => {
      return !resources.isLoading;
    },
    () => {
      resources.versionInfo('buildah');
      when(
        () => {
          return !resources.isLoading;
        },
        () => {
          setTimeout(() => {
            const component = mount(
              <Provider>
                <BasicDetails />
              </Provider>
            );
            component.update();

            const r = component.find(BasicDetails);
            expect(r.length).toEqual(1);

            expect(component.debug()).toMatchSnapshot();
            done();
          }, 1000);
        }
      );
    }
  );
});

it('length of DropdownItems should be 2 in case of buildah', (done) => {
  const { resources } = root;
  when(
    () => {
      return !resources.isLoading;
    },
    () => {
      resources.versionInfo('buildah');
      when(
        () => {
          return !resources.isLoading;
        },
        () => {
          setTimeout(() => {
            const component = mount(
              <Provider>
                <BasicDetails />
              </Provider>
            );
            component.update();

            const r = component.find(BasicDetails);
            expect(r.length).toEqual(1);

            expect(component.debug()).toMatchSnapshot();

            const c = component.find(Card);
            expect(c.find(CardActions).length).toBe(1);
            const dropdownItems = c.find(CardActions).find(Dropdown).props().dropdownItems;
            assert(dropdownItems);
            expect(dropdownItems.length).toBe(2);
            done();
          }, 1000);
        }
      );
    }
  );
});
