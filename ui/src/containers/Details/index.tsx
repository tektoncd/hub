import React from 'react';
import { useObserver } from 'mobx-react';
import { useParams } from 'react-router-dom';
import { Spinner } from '@patternfly/react-core';
import { useMst } from '../../store/root';

const Details: React.FC = () => {
  const { resources } = useMst();
  const { name } = useParams();

  return useObserver(() =>
    resources.isLoading ? (
      <Spinner className="hub-spinner" />
    ) : (
      <React.Fragment>
        <span>Add Resources detail here {resources.resources.get(name).name}</span>
      </React.Fragment>
    )
  );
};
export default Details;
