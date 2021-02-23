import { renderHook } from '@testing-library/react-hooks';
import { FakeHub } from '../api/testutil';
import { createProviderAndStore } from '../store/root';
import { useDebounce } from './useDebounce';

const TESTDATA_DIR = `src/store/testdata`;
const api = new FakeHub(TESTDATA_DIR);
const {} = createProviderAndStore(api);

describe('useDebounce', () => {
  it('it renders useDebounce', (done) => {
    const { result } = renderHook(() => useDebounce('cli', 400));
    expect(result.current).toBe('cli');

    done();
  });
});
