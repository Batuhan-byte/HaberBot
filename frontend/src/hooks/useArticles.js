import { useQuery } from '@tanstack/react-query';
import { fetchArticles } from '../services/api';
import { DEFAULT_PAGE, DEFAULT_LIMIT } from '../utils/constants';

/**
 * Custom hook to fetch paginated articles.
 * @param {number} page - Current page number
 * @param {number} limit - Items per page
 * @param {string} source - Optional source filter
 * @returns {object} TanStack Query result object
 */
export function useArticles(page = DEFAULT_PAGE, limit = DEFAULT_LIMIT, source = '') {
  return useQuery({
    queryKey: ['articles', page, limit, source],
    queryFn: () => fetchArticles(page, limit, source),
  });
}
