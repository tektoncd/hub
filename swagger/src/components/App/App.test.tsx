import { shallow } from 'enzyme';
import SwaggerUI from 'swagger-ui-react';
import App from '.';

test('should render the Swagger component', () => {
  const app = shallow(<App />);
  expect(SwaggerUI.length).toBe(1);
  expect(app.debug()).toMatchSnapshot();
});
