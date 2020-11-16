import React from 'react';
import { observer } from 'mobx-react';
import '@patternfly/react-core/dist/styles/base.css';
import { createProvider } from '../../store/root';
import LeftPane from '../../components/LeftPane';
import { Grid, GridItem } from '@patternfly/react-core';

const Provider = createProvider();

const App: React.FC = observer(() => {
  return (
    <Provider>
      <div className="App">
        <Grid>
          <GridItem span={2}>
            <LeftPane />
          </GridItem>
        </Grid>
      </div>
    </Provider>
  );
});

export default App;
