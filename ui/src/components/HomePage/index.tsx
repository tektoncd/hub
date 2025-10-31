import React from 'react';
import {
  Banner,
  Grid,
  GridItem,
  PageSection,
  TextContent,
  TextVariants,
  Text
} from '@patternfly/react-core';
import LeftPane from '../LeftPane';
import Background from '../Background';
import Resources from '../../containers/Resources';
import './HomePage.css';
const HomePage = () => {
  return (
    <>
      <Banner variant="warning" isSticky={true}>
        <TextContent>
          <Text component={TextVariants.h1} className="hub-banner-text">
            Tekton Hub is deprecated and will be shutdown on January 2026, the 7th.
          </Text>
        </TextContent>
      </Banner>
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
