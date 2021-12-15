import React from 'react';
import { mount } from 'enzyme';
import { when } from 'mobx';
import ReactMarkdown from 'react-markdown';
import { Card } from '@patternfly/react-core';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import Details from '.';
import BasicDetails from '../BasicDetails';
import Description from '../../components/Description';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider, root } = createProviderAndStore(api);

jest.mock('react-router-dom', () => {
  return {
    useHistory: () => {
      return {
        history: ''
      };
    },
    useParams: () => {
      return {
        name: 'buildah',
        kind: 'task',
        catalog: 'tekton',
        version: '0.1'
      };
    }
  };
});

describe('Details component', () => {
  it('should render the details component', (done) => {
    const component = mount(
      <Provider>
        <Details />
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

          const r = component.find(Details);
          expect(r.length).toEqual(1);

          expect(component.debug()).toMatchSnapshot();

          done();
        }, 1000);
      }
    );
  });

  it('should render the BasicDetails on details component', (done) => {
    const component = mount(
      <Provider>
        <Details />
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

          const r = component.find(Details);
          expect(r.length).toEqual(1);

          const c = component.find(BasicDetails);
          expect(c.length).toBe(1);
          expect(c.find(Card).length).toBe(1);
          done();
        }, 1000);
      }
    );
  });

  it('should render the Description on details component', (done) => {
    const component = mount(
      <Provider>
        <Details />
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

          const r = component.find(Details);
          expect(r.length).toEqual(1);

          expect(resources.filteredResources[0].name).toBe('buildah');

          const c = component.find(Description);
          expect(c.find(ReactMarkdown).length).toBe(2);

          done();
        }, 1000);
      }
    );
  });
});
