import React from 'react';
import { observer } from 'mobx-react';
import { useMst } from '../../store/root';
import Filter from '../../components/Filter';

const CategoryFilter: React.FC = observer(() => {
  const { categories } = useMst();
  return <Filter store={categories} header="Category" />;
});

export default CategoryFilter;
