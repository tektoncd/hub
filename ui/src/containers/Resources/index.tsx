import React from 'react';
import { useObserver } from 'mobx-react';
import {
  EmptyState,
  EmptyStateIcon,
  EmptyStateVariant,
  Gallery,
  Spinner,
  Title,
  Button
} from '@patternfly/react-core';
import CubesIcon from '@patternfly/react-icons/dist/js/icons/cubes-icon';
import { useHistory } from 'react-router-dom';
import { useMst } from '../../store/root';
import { IResource } from '../../store/resource';
import Cards from '../../components/Cards';
import './Resources.css';

const Resources = () => {
  const { resources } = useMst();

  const history = useHistory();
  const clearFilter = () => {
    resources.clearAllFilters();
    history.push('/');
  };

  const checkResources = (items: IResource[]) => {
    return !items.length ? (
      <EmptyState variant={EmptyStateVariant.full} className="hub-resource-emptystate__margin">
        <EmptyStateIcon icon={CubesIcon} />
        <Title headingLevel="h5" size="md">
          No Resource Found.
        </Title>
        <Button variant="primary" onClick={clearFilter}>
          Clear All Filters
        </Button>
      </EmptyState>
    ) : (
      <Gallery hasGutter className="hub-resource">
        <Cards items={items} />
      </Gallery>
    );
  };

  return useObserver(() =>
    resources.resources.size === 0 ? (
      <Spinner className="hub-spinner" />
    ) : (
      <React.Fragment>{checkResources(resources.filteredResources)}</React.Fragment>
    )
  );
};
export default Resources;
