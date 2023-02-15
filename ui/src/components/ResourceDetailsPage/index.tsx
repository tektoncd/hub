import React from 'react';
import { Grid, GridItem, PageSection } from '@patternfly/react-core';
import Details from '../../containers/Details';

const ResourceDetailsPage = () => {
  return (
    <PageSection>
      <Grid hasGutter>
        <GridItem span={12}>
          <Details />
        </GridItem>
      </Grid>
    </PageSection>
  );
};

export default ResourceDetailsPage;
