import React from 'react';
import { observer } from 'mobx-react';
import '@patternfly/react-core/dist/styles/base.css';

import CategoryFilter from '../CategoryFilter/CategoryFilter';
import { createProvider } from '../../store/root';

const Provider = createProvider();

const App: React.FC = observer(() => {
  return (
    <Provider>
      <div className="App">
        <CategoryFilter />
      </div>
    </Provider>
  );
});

export default App;
