import React from 'react';
import { useMst } from '../../store/root';
import { Params } from '../../common/params';

const ParseUrl: React.FC = () => {
  const { resources } = useMst();

  if (window.location.search) {
    const searchParams = new URLSearchParams(window.location.search);
    if (searchParams.has(Params.Query)) {
      resources.setSearch(searchParams.get(Params.Query));
    }
    if (searchParams.has(Params.SortBy)) {
      resources.setSortBy(searchParams.get(Params.SortBy));
    }

    // Storing url params to store inorder to parse the url only after successfully resource load
    resources.setURLParams(window.location.search);
  }

  return <> </>;
};
export default ParseUrl;
