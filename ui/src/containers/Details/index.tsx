import React from 'react';
import { observer } from 'mobx-react';
import { useParams } from 'react-router-dom';
import { Spinner } from '@patternfly/react-core';
import { useMst } from '../../store/root';
import BasicDetails from '../BasicDetails';
import Description from '../../components/Description';
import { assert } from '../../store/utils';
import { PageNotFound } from '../../components/PageNotFound';
import { titleCase } from '../../common/titlecase';
import { scrollToTop } from '../../common/scrollToTop';
import { IVersion } from '../../store/resource';

const Details: React.FC = observer(() => {
  const { resources, user } = useMst();
  const { name, catalog, kind, version } = useParams();

  React.useEffect(() => {
    if (resources.isResourceLoading === false) {
      const resourceKey = `${catalog}/${titleCase(kind as string)}/${name}`;
      resources.versionInfo(resourceKey);
    }
  }, [resources.isResourceLoading]);

  const resourceKey = `${catalog}/${titleCase(kind as string)}/${name}`;

  const validateUrl = () => {
    const resource = resources.resources.get(resourceKey);
    assert(resource);
    const allVersions = resource.versions;
    let isValidVersion = false;
    if (version !== undefined) {
      allVersions.forEach((item: IVersion) => {
        if (item.version === version) {
          isValidVersion = true;
        }
      });
    } else {
      isValidVersion = true;
    }
    return resources.resources.has(resourceKey) && isValidVersion;
  };

  const resourceDetails = () => {
    resources.loadReadme(resourceKey);
    resources.loadYaml(resourceKey);
    const resource = resources.resources.get(resourceKey);
    assert(resource);
    user.getRating(resource.id);
  };

  return resources.isVersionLoading === true ? (
    <Spinner className="hub-details-spinner" />
  ) : !validateUrl() ? (
    <PageNotFound />
  ) : (
    <React.Fragment>
      <>
        {resourceDetails()}
        {scrollToTop()}
        <BasicDetails />
        <Description name={name as string} catalog={catalog as string} kind={kind as string} />
      </>
    </React.Fragment>
  );
});
export default Details;
