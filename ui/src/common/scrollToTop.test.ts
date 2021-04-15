import { scrollToTop } from './scrollToTop';

describe('Test scrollToTop function', () => {
  it('should move the page to top', () => {
    scrollToTop();

    expect(window.screenX).toBe(0);
    expect(window.screenY).toBe(0);
  });
});
