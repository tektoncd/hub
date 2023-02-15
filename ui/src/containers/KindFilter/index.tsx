import React from 'react';
import { observer } from 'mobx-react';
import Filter from '../../components/Filter';
import { useMst } from '../../store/root';

const KindFilter: React.FC = observer(() => {
  const { resources } = useMst();
  return <Filter store={resources.kinds} header="Kind" />;
});

export default KindFilter;
