import React from 'react';
import { useObserver } from 'mobx-react';
import { GridItem, Grid } from '@patternfly/react-core';
import CatalogFilter from '../../containers/CatalogFilter';
import KindFilter from '../../containers/KindFilter';
import CategoryFilter from '../../containers/CategoryFilter';
import Sort from '../../containers/SortDropDown';
import './LeftPane.css';

const LeftPane: React.FC = () => {
  return useObserver(() => (
    <Grid hasGutter className="hub-leftpane">
      <GridItem span={8}>
        <Sort />
      </GridItem>
      <GridItem>
        <KindFilter />
      </GridItem>
      <GridItem>
        <CatalogFilter />
      </GridItem>
      <GridItem>
        <CategoryFilter />
      </GridItem>
    </Grid>
  ));
};
export default LeftPane;
