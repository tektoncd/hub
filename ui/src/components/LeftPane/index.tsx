import React from 'react';
import { useObserver } from 'mobx-react';
import { GridItem, Grid } from '@patternfly/react-core';
import CatalogFilter from '../../containers/CatalogFilter';
import KindFilter from '../../containers/KindFilter';
import CategoryFilter from '../../containers/CategoryFilter';
import PlatformFilter from '../../containers/PlatformFilter';
import Sort from '../../containers/SortDropDown';
import { useMst } from '../../store/root';
import { apiDownError } from '../../common/errors';
import './LeftPane.css';

const LeftPane: React.FC = () => {
  const { resources } = useMst();

  return useObserver(() =>
    resources.err !== apiDownError ? (
      <Grid hasGutter className="hub-leftpane">
        <GridItem span={8}>
          <Sort />
        </GridItem>
        <GridItem>
          <KindFilter />
        </GridItem>
        <GridItem>
          <PlatformFilter />
        </GridItem>
        <GridItem>
          <CatalogFilter />
        </GridItem>
        <GridItem>
          <CategoryFilter />
        </GridItem>
      </Grid>
    ) : null
  );
};
export default LeftPane;
