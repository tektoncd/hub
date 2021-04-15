import React from 'react';
import { useObserver } from 'mobx-react';
import { useParams } from 'react-router-dom';
import { Spinner } from '@patternfly/react-core';
import { useMst } from '../../store/root';
import BasicDetails from '../BasicDetails';
import Description from '../../components/Description';
import { assert } from '../../store/utils';
import { PageNotFound } from '../../components/PageNotFound';
import { titleCase } from '../../common/titlecase';
import { scrollToTop } from '../../common/scrollToTop';

const Details: React.FC = () => {
  const { resources, user } = useMst();
  const { name, catalog, kind } = useParams();

  const resourceKey = `${catalog}/${titleCase(kind)}/${name}`;
  const validateUrl = () => {
    return resources.resources.has(resourceKey);
  };

  const resourceDetails = () => {
    resources.versionInfo(resourceKey);
    resources.loadReadme(resourceKey);
    resources.loadYaml(resourceKey);
    const resource = resources.resources.get(resourceKey);
    assert(resource);
    user.getRating(resource.id);
  };

  return useObserver(() =>
    resources.resources.size === 0 ? (
      <Spinner className="hub-spinner" />
    ) : !validateUrl() ? (
      <PageNotFound />
    ) : (
      <React.Fragment>
        {resourceDetails()}
        {scrollToTop()}
        <BasicDetails />
        <Description name={name} catalog={catalog} kind={kind} />
      </React.Fragment>
    )
  );
};
export default Details;
