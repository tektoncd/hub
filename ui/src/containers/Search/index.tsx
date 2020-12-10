import { TextInput } from '@patternfly/react-core';
import React, { useState } from 'react';
import { useHistory } from 'react-router-dom';
import { useMst } from '../../store/root';
import './Search.css';

const Search: React.FC = () => {
  const { resources } = useMst();

  const [value, setValue] = useState('');

  const onSearchChange = (text: string) => {
    setValue(text);
    resources.setSearch(text.trim());
  };

  const history = useHistory();
  const onSearchKeyPress = (e: React.KeyboardEvent<HTMLElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      history.push('/');
    }
    return;
  };

  return (
    <TextInput
      value={value}
      type="search"
      onChange={onSearchChange}
      onKeyPress={onSearchKeyPress}
      aria-label="text input example"
      placeholder="Search for resources..."
      spellCheck="false"
      className="hub-search"
    />
  );
};

export default Search;
