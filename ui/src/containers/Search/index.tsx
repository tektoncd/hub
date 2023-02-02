import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useHotkeys } from 'react-hotkeys-hook';
import { useMst } from '../../store/root';
import { observer } from 'mobx-react';
import { TextInput } from '@patternfly/react-core';
import { useDebouncedCallback } from 'use-debounce';
import { assert } from '../../store/utils';
import './Search.css';

export const TagsKeyword = 'tags:';

const Search: React.FC = observer(() => {
  const { resources } = useMst();
  const [value, setValue] = React.useState(resources.search);

  interface SearchInfo {
    search: string;
    tags: Array<string>;
  }

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
    (data: SearchInfo) => {
      resources.setSearchedTags(data.tags);
      resources.setSearch(data.search);
    },
    // Delay in ms
    500
  );

  const onSearchChange = () => {
    const currentInput = inputRef.current;
    assert(currentInput);
    setValue(currentInput.value);

    const searchInput = currentInput.value;
    let tagList: Array<string>;

    if (searchInput.includes(TagsKeyword)) {
      const tagsString = searchInput.substr(5); // Removes `tags:` from search Input
      tagList = tagsString
        .split(',')
        .filter((item: string) => item !== '')
        .map((item: string) => item.trim());
    } else {
      tagList = [];
    }

    const searchedData: SearchInfo = {
      search: searchInput,
      tags: tagList
    };

    debounced(searchedData);
  };

  const history = useNavigate();
  const onSearchKeyPress = (e: React.KeyboardEvent<HTMLElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      if (window.location.pathname !== '/') history('/');
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
