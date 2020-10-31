import React from 'react';
import { useObserver } from 'mobx-react';
import Filter from '../Filter/Filter';
import { useMst } from '../../store/root';

const CategoryFilter: React.FC = () => {
  const { categories } = useMst();
  return useObserver(() => <Filter store={categories} header="Categories" />);
};

export default CategoryFilter;
