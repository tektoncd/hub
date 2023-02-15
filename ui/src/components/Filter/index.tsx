import React from 'react';
import { observer } from 'mobx-react';

import { Button, Checkbox, Text, TextVariants, Grid, GridItem } from '@patternfly/react-core';
import { IconSize, TimesIcon } from '@patternfly/react-icons';

import Icon from '../Icon';
import { titleCase } from '../../common/titlecase';

import './Filter.css';
import { Icons } from '../../common/icons';

interface Filterable {
  id?: number;
  name: string;
  selected: boolean;
  toggle(): void;
  icon?: Icons;
}

interface Store {
  values: Filterable[];
  clearSelected(): void;
}

interface Props {
  store: Store;
  header: string;
}

const labelWithIcon = (label: string, icon?: Icons) => (
  <div>
    {icon && <Icon id={icon} size={IconSize.sm} label={label} />}
    {titleCase(label)}
  </div>
);

const checkboxes = (values: Filterable[]) => {
  return values.map((c: Filterable) => (
    <Checkbox
      className="hub-filter-checkbox"
      key={c.name}
      label={labelWithIcon(c.name, c.icon)}
      isChecked={c.selected}
      onChange={() => c.toggle()}
      aria-label="controlled checkbox"
      id={`${c.name}`}
      name={c.name}
    />
  ));
};

const Filter: React.FC<Props> = observer(({ store, header }) => {
  return (
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
  );
});

export default Filter;
