import React, { ReactText } from 'react';
import { useObserver } from 'mobx-react';
import ReactMarkDown from 'react-markdown';
import gfm from 'remark-gfm';
import { Grid, Card, Tabs, Tab, GridItem, CardHeader, Spinner } from '@patternfly/react-core';
import { useMst } from '../../store/root';
import Readme from '../Readme';
import Yaml from '../Yaml';
import './Description.css';

interface Props {
  name: string;
}

const Description: React.FC<Props> = (props) => {
  const { resources } = useMst();

  const [activeTabKey, setActiveTabKey] = React.useState(0);
  const handleTabClick = (_: React.MouseEvent<HTMLElement, MouseEvent>, tabIndex: ReactText) => {
    setActiveTabKey(Number(tabIndex));
  };

  const resource = resources.resources.get(props.name);

  return useObserver(() =>
    resource.readme === '' || resource.yaml === '' ? (
      <Spinner className="hub-details-spinner" />
    ) : (
      <React.Fragment>
        <Grid className="hub-description">
          <GridItem offset={1} span={10}>
            <Card>
              <CardHeader className="hub-description-header">
                <Grid className="hub-tabs">
                  <GridItem span={12}>
                    <Tabs activeKey={activeTabKey} isSecondary onSelect={handleTabClick}>
                      <Tab eventKey={0} title="Description" id={props.name}>
                        <hr className="hub-horizontal-line"></hr>
                        <ReactMarkDown
                          plugins={[[gfm, { tablePipeAlign: false }]]}
                          source={resource.readme}
                          escapeHtml={true}
                          renderers={{ code: Readme }}
                          className="hub-readme"
                        />
                      </Tab>
                      <Tab eventKey={1} title="YAML" id={props.name}>
                        <hr className="hub-horizontal-line"></hr>
                        <ReactMarkDown
                          source={resource.yaml}
                          escapeHtml={true}
                          renderers={{ code: Yaml }}
                        />
                      </Tab>
                    </Tabs>
                  </GridItem>
                </Grid>
              </CardHeader>
            </Card>
          </GridItem>
        </Grid>
      </React.Fragment>
    )
  );
};

export default Description;
