import React from 'react';
import App from './App';
import { shallow } from 'enzyme';
import renderer from 'react-test-renderer';
import { FakeHub } from '../../api/testutil';
import { CategoryStore } from '../../store/category';
import CategoryFilter from '../CategoryFilter/CategoryFilter';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);

describe('App', () => {
  it('should render the component correctly and match the snapshot', (done) => {
    const store = CategoryStore.create({}, { api });
    const app = renderer.create(
      <div className="App">
        <CategoryFilter store={store} />
      </div>
    );
    expect(app.toJSON()).toMatchSnapshot();

    done();
  });

  it('should find the categoryFilter component and match the count', () => {
    const store = CategoryStore.create({}, { api });
    const component = shallow(<App store={store} />);

    expect(component.find(CategoryFilter).length).toEqual(1);
  });
});
