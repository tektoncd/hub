import React from 'react';
import { mount } from 'enzyme';
import { when } from 'mobx';
import { StarIcon } from '@patternfly/react-icons';
import Rating from '.';
import { FakeHub } from '../../api/testutil';
import { createProviderAndStore } from '../../store/root';

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
        name: 'buildah'
      };
    }
  };
});

describe('Rating ', () => {
  it('should render Ratings component', () => {
    const component = mount(
      <Provider>
        <Rating />
      </Provider>
    );
    expect(component.debug()).toMatchSnapshot();
  });

  it('should find star Icon', () => {
    const component = mount(
      <Provider>
        <Rating />
      </Provider>
    );
    expect(component.find(StarIcon).length).toBe(5);
  });

  it('should update user rating for a resource', (done) => {
    const { user } = root;

    const code = {
      code: 'foo'
    };

    user.authenticate(code, 'baar');

    when(
      () => !user.isLoading,
      () => {
        user.getRating(13);
        when(
          () => !user.isLoading,
          () => {
            expect(user.userRating).toBe(2);
            user.setRating(13, 3);

            when(
              () => !user.isLoading,
              () => {
                const component = mount(
                  <Provider>
                    <Rating />
                  </Provider>
                );

                expect(component.find('input').get(Number('2')).props.checked).toBe(true);
                done();
              }
            );
          }
        );
      }
    );
  });

  it('should find user rating for a resource', (done) => {
    const { user } = root;
    const code = {
      code: 'foo'
    };
    user.authenticate(code, 'baar');
    when(
      () => !user.isLoading,
      () => {
        user.getRating(13);
        when(
          () => !user.isLoading,
          () => {
            expect(user.userRating).toBe(2);
            const component = mount(
              <Provider>
                <Rating />
              </Provider>
            );
            expect(component.find('input').get(Number('3')).props.checked).toBe(true);
            done();
          }
        );
      }
    );
  });
});
