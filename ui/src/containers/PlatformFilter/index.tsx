import React from 'react';
import { observer } from 'mobx-react';
import Filter from '../../components/Filter';
import { useMst } from '../../store/root';

const PlatformFilter: React.FC = observer(() => {
  const { resources } = useMst();
  return <Filter store={resources.platforms} header="Platform" />;
});

export default PlatformFilter;
