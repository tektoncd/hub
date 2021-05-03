import React from 'react';
import { useHistory } from 'react-router-dom';
import { useObserver } from 'mobx-react';
import { useHotkeys } from 'react-hotkeys-hook';
import { TextInput } from '@patternfly/react-core';
import { useMst } from '../../store/root';
import { useDebounce } from '../../utils/useDebounce';
import { assert } from '../../store/utils';
import './Search.css';

const Search: React.FC = () => {
  useHotkeys('/', () => {
    // Without timeout the hotkey will be sent to select input
    setTimeout(() => {
      const currentInput = inputRef.current;
      assert(currentInput);
      currentInput.focus();
    }, 20);
  });

  const { resources } = useMst();

  const inputRef = React.useRef<HTMLInputElement>(null);

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
      ref={inputRef}
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
