import React, { ReactText } from 'react';
import { useObserver } from 'mobx-react';
import ReactMarkDown from 'react-markdown';
import gfm from 'remark-gfm';
import { Grid, Card, Tabs, Tab, GridItem, CardHeader, Spinner } from '@patternfly/react-core';
import { useMst } from '../../store/root';
import Readme from '../Readme';
import Yaml from '../Yaml';
import { titleCase } from '../../common/titlecase';
import './Description.css';
import { assert } from '../../store/utils';

interface Props {
  name: string;
  catalog: string;
  kind: string;
}

const Description: React.FC<Props> = (props) => {
  const { resources } = useMst();

  const [activeTabKey, setActiveTabKey] = React.useState(0);
  const handleTabClick = (_: React.MouseEvent<HTMLElement, MouseEvent>, tabIndex: ReactText) => {
    setActiveTabKey(Number(tabIndex));
  };

  const { catalog, kind, name } = props;
  const resource = resources.resources.get(`${catalog}/${titleCase(kind)}/${name}`);
  assert(resource);

  const { webURL, version } = resource.displayVersion;

  const resourceDirUrl = webURL.slice(0, webURL.indexOf(name));
  const resourceWebUrl = resourceDirUrl + `${name}`;

  // This function transform relative uri of readme into absoulte uri
  const transformUri = (uri: string) => {
    if (!uri.includes('./') && !uri.includes('http')) {
      return resourceWebUrl + `/${version}/${uri}`;
    }

    if (uri.includes('./')) {
      const uriPath = uri.slice(uri.lastIndexOf('./') + 1);

      if (uri.includes('../../')) {
        return resourceDirUrl + uriPath;
      }

      if (!/\d/.test(uriPath.slice(0, 3))) {
        return resourceWebUrl + `/${version}${uriPath}`;
      }

      return resourceWebUrl + uriPath;
    }

    return uri;
  };

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
                          linkTarget={' '}
                          transformLinkUri={(uri: string) => transformUri(uri)}
                          remarkPlugins={[[gfm, { tablePipeAlign: false }]]}
                          skipHtml={true}
                          components={{
                            code: ({ ...props }) => <Readme value={props.children as string} />
                          }}
                          className="hub-readme"
                        >
                          {resource.readme}
                        </ReactMarkDown>
                      </Tab>
                      <Tab eventKey={1} title="YAML" id={props.name}>
                        <hr className="hub-horizontal-line"></hr>
                        <ReactMarkDown
                          skipHtml={true}
                          components={{
                            code: ({ ...props }) => <Yaml value={props.children as string} />,
                            pre: ({ ...props }) => <span style={{ color: 'green' }} {...props} />
                          }}
                        >
                          {resource.yaml}
                        </ReactMarkDown>
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
