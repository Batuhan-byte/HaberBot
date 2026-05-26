import { useQuery } from '@tanstack/react-query';
import { fetchTopics } from '../services/api';

/**
 * Custom hook to fetch all available topics.
 * @returns {object} TanStack Query result object
 */
export function useTopics() {
  return useQuery({
    queryKey: ['topics'],
    queryFn: fetchTopics,
  });
}
