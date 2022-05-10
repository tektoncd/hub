import React from 'react';
import { observer } from 'mobx-react';
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
import { UpdateURL } from '../../utils/updateUrl';
import './Resources.css';

const Resources: React.FC = observer(() => {
  const { resources, categories } = useMst();
  const { catalogs, kinds, platforms, search, sortBy, searchedTags } = resources;

  const history = useHistory();

  React.useEffect(() => {
    const selectedcategories = categories.selectedByName.join(',');
    const selectedKinds = Array.from(kinds.selected).join(',');
    const selectedCatalogs = catalogs.selectedByName.join(',');
    const selectedPlatforms = platforms.selectedByName.join(',');

    const url = UpdateURL(
      search,
      sortBy,
      selectedcategories,
      selectedKinds,
      selectedCatalogs,
      selectedPlatforms,
      searchedTags
    );
    if (!resources.isLoading) history.replace(`?${url}`);
  }, [
    search,
    sortBy,
    categories.selectedByName,
    kinds.selected,
    catalogs.selected,
    platforms.selectedByName,
    searchedTags
  ]);

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

  return resources.resources.size === 0 ? (
    <Spinner className="hub-resources-spinner" />
  ) : (
    <React.Fragment>{checkResources(resources.filteredResources)}</React.Fragment>
  );
});
export default Resources;
