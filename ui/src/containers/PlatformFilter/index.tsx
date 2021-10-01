import React from 'react';
import { useObserver } from 'mobx-react';
import Filter from '../../components/Filter';
import { useMst } from '../../store/root';

const PlatformFilter: React.FC = () => {
  const { resources } = useMst();
  return useObserver(() => <Filter store={resources.platforms} header="Platform" />);
};

export default PlatformFilter;
