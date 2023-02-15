import React from 'react';
import { Grid, GridItem, PageSection } from '@patternfly/react-core';
import LeftPane from '../LeftPane';
import Background from '../Background';
import Resources from '../../containers/Resources';

const HomePage = () => {
  return (
    <>
      <Background />
      <PageSection>
        <Grid hasGutter>
          <GridItem span={2}>
            <LeftPane />
          </GridItem>
          <GridItem span={10} rowSpan={1}>
            <Resources />
          </GridItem>
        </Grid>
      </PageSection>
    </>
  );
};

export default HomePage;
