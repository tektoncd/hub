import React from 'react';
import { useHistory } from 'react-router-dom';
import { useHotkeys } from 'react-hotkeys-hook';
import { useMst } from '../../store/root';
import { observer } from 'mobx-react';
import { TextInput } from '@patternfly/react-core';
import { useDebouncedCallback } from 'use-debounce';
import { assert } from '../../store/utils';
import './Search.css';

const Search: React.FC = observer(() => {
  const { resources } = useMst();
  const [value, setValue] = React.useState(resources.search);

  React.useEffect(() => {
    setValue(resources.search);
  }, [resources.search]);

  const inputRef = React.useRef<HTMLInputElement>(null);
  useHotkeys('/', () => {
    // Without timeout the hotkey will be sent to select input
    setTimeout(() => {
      assert(inputRef.current);
      inputRef.current.focus();
    }, 20);
  });

  /* Debounce callback
     This function delay the invoking func for 500ms after
     when users stop searching
  */
  const debounced = useDebouncedCallback(
    (value: string) => {
      resources.setSearch(value);
    },
    // Delay in ms
    500
  );

  const onSearchChange = () => {
    const currentInput = inputRef.current;
    assert(currentInput);
    setValue(currentInput.value);
    debounced(currentInput.value);
  };

  const history = useHistory();
  const onSearchKeyPress = (e: React.KeyboardEvent<HTMLElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      if (window.location.pathname !== '/') history.push('/');
    }
    return;
  };

  return (
    <TextInput
      ref={inputRef}
      type="search"
      value={value}
      onChange={onSearchChange}
      onKeyPress={onSearchKeyPress}
      aria-label="text input example"
      placeholder="Search for resources..."
      spellCheck="false"
      className="hub-search"
    />
  );
});

export default Search;
