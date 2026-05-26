import { useQuery } from '@tanstack/react-query';
import { fetchTopicArticles } from '../services/api';
import { DEFAULT_PAGE, DEFAULT_LIMIT } from '../utils/constants';

/**
 * Custom hook to fetch articles for a specific topic.
 * @param {string} slug - Topic slug
 * @param {number} page - Current page number
 * @param {number} limit - Items per page
 * @returns {object} TanStack Query result object
 */
export function useTopicArticles(slug, page = DEFAULT_PAGE, limit = DEFAULT_LIMIT) {
  return useQuery({
    queryKey: ['topicArticles', slug, page, limit],
    queryFn: () => fetchTopicArticles(slug, page, limit),
    enabled: !!slug,
  });
}
