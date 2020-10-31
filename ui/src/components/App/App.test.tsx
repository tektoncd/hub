import React from 'react';
import App from '.';
import { shallow } from 'enzyme';
import renderer from 'react-test-renderer';
import { FakeHub } from '../../api/testutil';
import CategoryFilter from '../CategoryFilter/CategoryFilter';
import { createProviderAndStore } from '../../store/root';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { Provider } = createProviderAndStore(api);

describe('App', () => {
  it('should render the component correctly and match the snapshot', (done) => {
    const app = renderer.create(
      <Provider>
        <div className="App">
          <CategoryFilter />
        </div>
      </Provider>
    );

    expect(app.toJSON()).toMatchSnapshot();
    done();
  });

  it('should find the categoryFilter component and match the count', () => {
    const component = shallow(<App />);
    expect(component.find(CategoryFilter).length).toEqual(1);
  });
});
