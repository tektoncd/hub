import React from 'react';
import { useObserver } from 'mobx-react';
import Filter from '../Filter/Filter';
import { ICategoryStore } from '../../store/category';

interface store {
  store: ICategoryStore;
}

const CategoryFilter: React.FC<store> = (props: store) => {
  const store = props.store;
  return useObserver(() => <Filter store={store} header="Categories" />);
};

export default CategoryFilter;
