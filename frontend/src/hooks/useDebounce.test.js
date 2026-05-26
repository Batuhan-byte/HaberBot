import { renderHook, act } from '@testing-library/react';
import { useDebounce } from './useDebounce';

describe('useDebounce', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('returns initial value immediately', () => {
    const { result } = renderHook(() => useDebounce('hello', 300));
    expect(result.current).toBe('hello');
  });

  it('updates debounced value after delay', () => {
    const { result, rerender } = renderHook(
      ({ value, delay }) => useDebounce(value, delay),
      { initialProps: { value: 'hello', delay: 300 } }
    );

    rerender({ value: 'hello world', delay: 300 });

    expect(result.current).toBe('hello');

    act(() => { vi.advanceTimersByTime(300); });

    expect(result.current).toBe('hello world');
  });

  it('resets timer when value changes before delay', () => {
    const { result, rerender } = renderHook(
      ({ value, delay }) => useDebounce(value, delay),
      { initialProps: { value: 'a', delay: 300 } }
    );

    rerender({ value: 'ab', delay: 300 });
    act(() => { vi.advanceTimersByTime(200); });

    rerender({ value: 'abc', delay: 300 });
    act(() => { vi.advanceTimersByTime(200); });

    expect(result.current).toBe('a');

    act(() => { vi.advanceTimersByTime(100); });

    expect(result.current).toBe('abc');
  });

  it('uses default delay of 300ms', () => {
    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value),
      { initialProps: { value: 'first' } }
    );

    rerender({ value: 'second' });
    act(() => { vi.advanceTimersByTime(299); });
    expect(result.current).toBe('first');

    act(() => { vi.advanceTimersByTime(1); });
    expect(result.current).toBe('second');
  });
});
