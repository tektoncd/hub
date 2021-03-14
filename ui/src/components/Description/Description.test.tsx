import React from 'react';
import { mount } from 'enzyme';
import { when } from 'mobx';
import ReactMarkDown from 'react-markdown';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';
import Description from '.';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider, root } = createProviderAndStore(api);

describe('Resource Readme and Yaml', () => {
  it('render readme and yaml', (done) => {
    const { resources } = root;
    when(
      () => {
        return !resources.isLoading;
      },
      () => {
        resources.loadReadme('tekton/Task/buildah');
        resources.loadYaml('tekton/Task/buildah');
        when(
          () => {
            return !resources.isLoading;
          },
          () => {
            setTimeout(() => {
              const component = mount(
                <Provider>
                  <Description name="buildah" catalog="tekton" kind="Task" />
                </Provider>
              );
              component.update();

              expect(component.find(ReactMarkDown).length).toBe(2);

              done();
            }, 1000);
          }
        );
      }
    );
  });
});
