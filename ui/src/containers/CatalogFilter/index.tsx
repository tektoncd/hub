import React from 'react';
import { observer } from 'mobx-react';
import Filter from '../../components/Filter';
import { useMst } from '../../store/root';

const CatalogFilter: React.FC = observer(() => {
  const { catalogs } = useMst();
  return <Filter store={catalogs} header="Catalog" />;
});

export default CatalogFilter;
