import { useQuery } from '@tanstack/react-query';
import { fetchArticle } from '../services/api';

/**
 * Custom hook to fetch a single article by ID.
 * @param {string|number} articleId - The article ID
 * @returns {object} TanStack Query result object
 */
export function useArticle(articleId) {
  return useQuery({
    queryKey: ['article', articleId],
    queryFn: () => fetchArticle(articleId),
    enabled: !!articleId,
  });
}
