import {
  Banner,
  Grid,
  GridItem,
  PageSection,
  Text,
  TextContent,
  TextVariants
} from '@patternfly/react-core';
import Resources from '../../containers/Resources';
import Background from '../Background';
import LeftPane from '../LeftPane';
import './HomePage.css';
const HomePage = () => {
  return (
    <>
      <Banner variant="warning" isSticky={true}>
        <TextContent>
          <Text component={TextVariants.h1} className="hub-banner-text">
            Tekton Hub will be deprecated in v1.24.0 and removed in v1.26.0.
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
