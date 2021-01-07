import React from 'react';
import { useObserver } from 'mobx-react';
import {
  EmptyState,
  EmptyStateIcon,
  EmptyStateVariant,
  Gallery,
  Spinner,
  Title
} from '@patternfly/react-core';
import CubesIcon from '@patternfly/react-icons/dist/js/icons/cubes-icon';
import { useMst } from '../../store/root';
import { IResource } from '../../store/resource';
import Cards from '../../components/Cards';
import './Resources.css';

const Resources = () => {
  const { resources } = useMst();

  const checkResources = (resources: IResource[]) => {
    return !resources.length ? (
      <EmptyState variant={EmptyStateVariant.full} className="hub-resource-emptystate__margin">
        <EmptyStateIcon icon={CubesIcon} />
        <Title headingLevel="h5" size="md">
          No Resource Found.
        </Title>
      </EmptyState>
    ) : (
      <Gallery hasGutter className="hub-resource">
        <Cards items={resources} />
      </Gallery>
    );
  };

  return useObserver(() =>
    resources.isLoading ? (
      <Spinner className="hub-spinner" />
    ) : (
      <React.Fragment>{checkResources(resources.filteredResources)}</React.Fragment>
    )
  );
};
export default Resources;
