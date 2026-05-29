import { useQuery } from '@tanstack/react-query';
import { api } from '../services/api';

/**
 * Custom hook to fetch the AI summary of a single article by ID.
 * Caches the response indefinitely (client session cache) by setting staleTime and gcTime/cacheTime.
 * @param {string|number} articleId - The article ID
 * @param {boolean} enabled - Whether to fetch immediately or wait
 * @returns {object} TanStack Query result object
 */
export function useSummary(articleId, enabled = false) {
  return useQuery({
    queryKey: ['article-summary', articleId],
    queryFn: () => api.getArticleSummary(articleId),
    enabled: enabled && !!articleId,
    staleTime: Infinity,
    gcTime: Infinity,
  });
}
