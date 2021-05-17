import React from 'react';
import { Card, Grid, Text, GridItem, TextVariants, TextContent } from '@patternfly/react-core';
import '@patternfly/react-core/dist/styles/base.css';
import tekton from '../../assets/logo/tekton.png';
import './Footer.css';

const Footer: React.FC = () => {
  return (
    <React.Fragment>
      <Card className="hub-footer-card">
        <Grid>
          <GridItem span={12} className="hub-footer-info">
            <a href="https://cd.foundation" target="_">
              <img
                src={`https://tekton.dev/partner-logos/cdf.png?${global.Date.now()}`}
                alt="tekton.dev"
              />
            </a>
          </GridItem>
          <GridItem span={12} className="hub-footer-info">
            <TextContent className="hub-info-color">
              <Text component={TextVariants.h3}>
                Tekton is a{' '}
                <Text component={TextVariants.a} href="https://cd.foundation" target="_">
                  Continuous Delivery Foundation
                </Text>{' '}
                project.
              </Text>
            </TextContent>
          </GridItem>
          <GridItem span={12} className="hub-logo-margin">
            <img src={`${tekton}?${global.Date.now()}`} alt="Tekton" className="hub-logo-size" />
          </GridItem>
          <GridItem span={12} className="hub-footer-description">
            <Text>
              © {new Date().getFullYear()} The Linux Foundation®. All rights reserved. The Linux
              Foundation has registered trademarks and uses trademarks. For a list of trademarks of
              The Linux Foundation, please see our{' '}
              <Text
                component={TextVariants.a}
                href="https://www.linuxfoundation.org/trademark-usage/"
                target="_"
              >
                Trademark Usage page
              </Text>
              . Linux is a registered trademark of Linus Torvalds.{' '}
              <Text
                component={TextVariants.a}
                href="https://www.linuxfoundation.org/privacy/"
                target="_"
              >
                Privacy Policy
              </Text>{' '}
              and{' '}
              <Text
                component={TextVariants.a}
                href="https://www.linuxfoundation.org/terms/"
                target="_"
              >
                Terms of Use
              </Text>{' '}
              .
            </Text>
          </GridItem>
        </Grid>
      </Card>
    </React.Fragment>
  );
};
export default Footer;
