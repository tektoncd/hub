import React, { useEffect } from 'react';
import { useHistory } from 'react-router-dom';
import { useObserver } from 'mobx-react';
import { TextInput } from '@patternfly/react-core';
import { useMst } from '../../store/root';
import { useDebounce } from '../../utils/useDebounce';
import './Search.css';

const Search: React.FC = () => {
  const { resources } = useMst();

  // to get query params from the url
  const searchParams = new URLSearchParams(window.location.search);
  const query = searchParams.get('query') || ' ';

  useEffect(() => {
    if (query !== ' ') {
      resources.setSearch(query);
    }
  }, [query, resources]);

  const setParams = ({ query = '' }) => {
    const searchParams = new URLSearchParams();
    searchParams.set('query', query);
    return searchParams.toString();
  };

  const updateURL = (text: string) => {
    const url = setParams({ query: text });
    if (window.location.pathname === '/') history.replace(`?${url}`);
  };

  const onSearchChange = useDebounce(resources.search, 400);

  const history = useHistory();
  const onSearchKeyPress = (e: React.KeyboardEvent<HTMLElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      history.push('/');
      updateURL(resources.search);
    }
    return;
  };

  return useObserver(() => (
    <TextInput
      value={resources.search}
      type="search"
      onChange={(resourceName: string) => {
        resources.setSearch(resourceName);
        updateURL(resourceName);
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
