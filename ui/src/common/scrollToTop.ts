import { assert } from '../store/utils';

// This function moves the page to top
export const scrollToTop = () => {
  const scroller = document.querySelector('main');
  assert(scroller);
  if (scroller) scroller.scrollTo(0, 0);
};
