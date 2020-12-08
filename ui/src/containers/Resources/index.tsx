import React from 'react';
import { Gallery, Spinner } from '@patternfly/react-core';
import { useMst } from '../../store/root';
import { useObserver } from 'mobx-react';
import Cards from '../../components/Cards';
import './Resources.css';

const Resources = () => {
  const { resources } = useMst();

  return useObserver(() =>
    resources.isLoading ? (
      <Spinner className="hub-spinner" />
    ) : (
      <React.Fragment>
        <Gallery hasGutter className="hub-resource">
          <Cards items={resources.filteredResources} />
        </Gallery>
      </React.Fragment>
    )
  );
};

export default Resources;
