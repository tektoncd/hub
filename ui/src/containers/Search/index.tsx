import { TextInput } from '@patternfly/react-core';
import React, { useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import { useMst } from '../../store/root';
import { useDebounce } from '../../utils/useDebounce';
import './Search.css';

const Search: React.FC = () => {
  const { resources } = useMst();
  const [value, setValue] = useState('');

  // to get query params from the url
  const searchParams = new URLSearchParams(window.location.search);
  const query = searchParams.get('query') || ' ';

  useEffect(() => {
    if (query !== ' ') {
      resources.setSearch(query);
      setValue(query);
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

  const onSearchChange = useDebounce(value, 400);

  const history = useHistory();
  const onSearchKeyPress = (e: React.KeyboardEvent<HTMLElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      history.push('/');
      updateURL(value);
    }
    return;
  };

  return (
    <TextInput
      value={value}
      type="search"
      onChange={(resourceName: string) => {
        setValue(resourceName);
        updateURL(resourceName);
        return onSearchChange;
      }}
      onKeyPress={onSearchKeyPress}
      aria-label="text input example"
      placeholder="Search for resources..."
      spellCheck="false"
      className="hub-search"
    />
  );
};

export default Search;
