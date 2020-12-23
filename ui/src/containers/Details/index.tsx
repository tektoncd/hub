import React from 'react';
import { useObserver } from 'mobx-react';
import { useParams } from 'react-router-dom';
import { Spinner } from '@patternfly/react-core';
import { useMst } from '../../store/root';
import BasicDetails from '../BasicDetails';
import Description from '../../components/Description';
import { assert } from '../../store/utils';

const Details: React.FC = () => {
  const { resources, user } = useMst();
  const { name } = useParams();

  const resourceDetails = () => {
    resources.versionInfo(name);
    resources.loadReadme(name);
    resources.loadYaml(name);

    const resource = resources.resources.get(name);
    assert(resource);
    user.getRating(resource.id);
  };

  return useObserver(() =>
    resources.resources.size === 0 ? (
      <Spinner className="hub-spinner" />
    ) : (
      <React.Fragment>
        {resourceDetails()}
        <BasicDetails />
        <Description name={name} />
      </React.Fragment>
    )
  );
};
export default Details;
