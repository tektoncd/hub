import React from 'react';
import {
  EmptyState,
  EmptyStateVariant,
  Title,
  Button,
  EmptyStateIcon
} from '@patternfly/react-core';
import { IconSize } from '@patternfly/react-icons';
import { useHistory } from 'react-router-dom';
import { Icons } from '../../common/icons';
import Icon from '../Icon';
import './PageNotFound.css';

export const PageNotFound: React.FC = () => {
  const history = useHistory();

  const warningIcon = () => {
    return <Icon id={Icons.WarningTriangle} size={IconSize.xl} label="Page Not Found" />;
  };

  return (
    <EmptyState variant={EmptyStateVariant.small}>
      <EmptyStateIcon icon={warningIcon} />
      <Title headingLevel="h5" size="lg" className="hub-pagenotfound__title">
        Page Not Found
      </Title>
      <Button variant="primary" onClick={() => history.push('/')}>
        Back Home
      </Button>
    </EmptyState>
  );
};
