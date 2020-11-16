import React, { useState } from 'react';
import { DropdownItem, Dropdown, DropdownToggle } from '@patternfly/react-core';
import { useObserver } from 'mobx-react';
import './SortDropDown.css';
import { useMst } from '../../store/root';
import { sortByFields } from '../../store/resource';

const Sort: React.FC = () => {
  const { resources } = useMst();

  const items: Array<string> = Object.values(sortByFields);

  const [isOpen, set] = useState(false);

  const dropDownItems = items.map((value) => (
    <DropdownItem id={value} key={value} onClick={(e) => resources.setSortBy(e.currentTarget.id)}>
      {value}
    </DropdownItem>
  ));

  const onToggle = (isOpen: React.SetStateAction<boolean>) => set(isOpen);
  const onSelect = () => set(!isOpen);

  return useObserver(() => {
    return (
      <div>
        <Dropdown
          onSelect={onSelect}
          toggle={<DropdownToggle onToggle={onToggle}>{resources.sortBy}</DropdownToggle>}
          isOpen={isOpen}
          dropdownItems={dropDownItems}
          className="hub-sort-dropdown"
        />
      </div>
    );
  });
};

export default Sort;
