import React from 'react';
import { observer } from 'mobx-react';
import '@patternfly/react-core/dist/styles/base.css';
import { Grid, GridItem, Page, PageSection } from '@patternfly/react-core';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import LeftPane from '../../components/LeftPane';
import Background from '../../components/Background';
import Header from '../../containers/Header';
import Footer from '../../components/Footer';
import Resources from '../Resources';
import Details from '../Details';
import ParseUrl from '../ParseUrl';
import { createProvider } from '../../store/root';
import './App.css';

const Provider = createProvider();

const App: React.FC = observer(() => {
  return (
    <Provider>
      <Router>
        <ParseUrl />
        <Page header={<Header />} className="hub-page">
          <Route exact path="/" component={Background} />
          <PageSection>
            <Grid hasGutter>
              <GridItem span={12}>
                <Route exact path="/:catalog/:kind/:name/:version?" component={Details}></Route>
              </GridItem>
              <GridItem span={2}>
                <Route exact path="/" component={LeftPane}></Route>
              </GridItem>

              <GridItem span={10} rowSpan={1}>
                <Route exact path="/" component={Resources}></Route>
              </GridItem>
            </Grid>
          </PageSection>
          <Footer />
        </Page>
      </Router>
    </Provider>
  );
});

export default App;
