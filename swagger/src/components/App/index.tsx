import React from 'react';
import SwaggerUI from 'swagger-ui-react';
import 'swagger-ui-react/swagger-ui.css';
import { API_URL } from '../../config/constants';

const App: React.FC = () => {
  return (
    <div>
      <SwaggerUI url={API_URL} />
    </div>
  );
};

export default App;
