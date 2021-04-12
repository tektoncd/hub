import React from 'react';
import { useObserver } from 'mobx-react';
import Filter from '../../components/Filter';
import { useMst } from '../../store/root';

const CatalogFilter: React.FC = () => {
  const { catalogs } = useMst();
  return useObserver(() => <Filter store={catalogs} header="Catalog" />);
};

export default CatalogFilter;
