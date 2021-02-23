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
import { ICatalog } from '../../store/catalog';

const Details: React.FC = () => {
  const { resources, user } = useMst();
  const { name, catalog, kind } = useParams();

  const catalogs = resources.catalogs.values;
  const validateUrl = () => {
    let catalogUrl = false;
    catalogs.forEach((item: ICatalog) => {
      if (item.name === catalog) catalogUrl = true;
    });
    return (
      catalogUrl && resources.kinds.items.has(titleCase(kind)) && resources.resources.has(name)
    );
  };

  const resourceDetails = () => {
    resources.versionInfo(name);
    resources.loadReadme(name);
    resources.loadYaml(name);
    const resource = resources.resources.get(name);
    assert(resource);
    user.getRating(resource.id);
  };

  const scrollToTop = () => {
    const scroller = document.querySelector('main');
    assert(scroller);
    if (scroller) scroller.scrollTo(0, 0);
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
        <Description name={name} />
      </React.Fragment>
    )
  );
};
export default Details;
