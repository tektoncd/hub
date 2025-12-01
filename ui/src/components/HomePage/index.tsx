import { Grid, GridItem, PageSection } from '@patternfly/react-core';
import Resources from '../../containers/Resources';
import Background from '../Background';
import LeftPane from '../LeftPane';
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
