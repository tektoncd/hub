import React from 'react';
import { Text, TextContent, TextVariants } from '@patternfly/react-core';
import './Background.css';

const Background: React.FC = () => (
  <div className="hub-background-image">
    <TextContent>
      <Text component={TextVariants.h1} className="hub-background-heading">
        Welcome to Tekton Hub
      </Text>
      <Text component={TextVariants.h3} className="hub-tag-line">
        Discover, search and share reusable Tasks and Pipelines
      </Text>
    </TextContent>
  </div>
);

export default Background;
