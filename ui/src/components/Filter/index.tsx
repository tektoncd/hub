import React from 'react';
import { useObserver } from 'mobx-react';
import { Button, Checkbox, Text, TextVariants, Grid, GridItem } from '@patternfly/react-core';
import TimesIcon from '@patternfly/react-icons/dist/js/icons/times-icon';
import './Filter.css';

interface Filterable {
  id?: number;
  name: string;
  selected: boolean;
  toggle(): void;
}

interface Store {
  values: Filterable[];
  clearSelected(): void;
}

interface Props {
  store: Store;
  header: string;
}

const checkboxes = (values: Filterable[]) => {
  return values.map((c: Filterable) => (
    <Checkbox
      key={c.name}
      label={c.name}
      isChecked={c.selected}
      onChange={() => c.toggle()}
      aria-label="controlled checkbox"
      id={`${c.id}`}
      name={c.name}
    />
  ));
};

const Filter: React.FC<Props> = ({ store, header }) => {
  return useObserver(() => (
    <div>
      <Grid>
        <GridItem span={6}>
          <Text component={TextVariants.h1} className="hub-filter-header">
            {header}
          </Text>
        </GridItem>
        <GridItem span={3}>
          <Button variant="plain" aria-label="Clear" onClick={store.clearSelected}>
            <TimesIcon />
          </Button>
        </GridItem>

        <GridItem span={12}>{checkboxes(store.values)}</GridItem>
      </Grid>
    </div>
  ));
};

export default Filter;
