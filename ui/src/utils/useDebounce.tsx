import { useEffect, useState } from 'react';
import { useMst } from '../store/root';

export const useDebounce = (value: string, delay: number) => {
  const [debouncedValue, setDebouncedValue] = useState(value);
  const { resources } = useMst();

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
      resources.setSearch(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay, resources]);
  return debouncedValue;
};
