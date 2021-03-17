import React from 'react';
import { useHistory } from 'react-router-dom';
import { useObserver } from 'mobx-react';
import { TextInput } from '@patternfly/react-core';
import { useMst } from '../../store/root';
import { useDebounce } from '../../utils/useDebounce';
import './Search.css';

const Search: React.FC = () => {
  const { resources } = useMst();

  const onSearchChange = useDebounce(resources.search, 400);

  const history = useHistory();
  const onSearchKeyPress = (e: React.KeyboardEvent<HTMLElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      if (window.location.pathname !== '/') history.push('/');
    }
    return;
  };

  return useObserver(() => (
    <TextInput
      value={resources.search}
      type="search"
      onChange={(resourceName: string) => {
        resources.setSearch(resourceName);
        return onSearchChange;
      }}
      onKeyPress={onSearchKeyPress}
      aria-label="text input example"
      placeholder="Search for resources..."
      spellCheck="false"
      className="hub-search"
    />
  ));
};

export default Search;
