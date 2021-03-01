import { renderHook } from '@testing-library/react-hooks';
import { when } from 'mobx';
import { FakeHub } from '../api/testutil';
import { createProviderAndStore } from '../store/root';
import { useDebounce } from './useDebounce';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const { root } = createProviderAndStore(api);

describe('useDebounce', () => {
  it('it renders useDebounce', (done) => {
    const { resources } = root;
    when(
      () => {
        return !resources.isLoading;
      },
      () => {
        setTimeout(() => {
          const { result } = renderHook(() => useDebounce('cli', 400));
          expect(result.current).toBe('cli');
          done();
        }, 0);
      }
    );
  });
});
