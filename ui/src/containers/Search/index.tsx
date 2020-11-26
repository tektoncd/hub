import { TextInput } from '@patternfly/react-core';
import React, { useState } from 'react';
import './Search.css';

const Search: React.FC = () => {
  const [value, setValue] = useState('');

  const setInput = (event: string) => {
    setValue(event);
  };

  return (
    <TextInput
      value={value}
      type="search"
      onChange={setInput}
      aria-label="text input example"
      placeholder="Search for resources..."
      spellCheck="false"
      className="hub-search"
    />
  );
};

export default Search;
