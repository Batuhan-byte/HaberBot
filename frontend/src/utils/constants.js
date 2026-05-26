/** API base URL from environment variable */
export const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

/** API version prefix */
export const API_PREFIX = '/api/v1';

/** Default pagination */
export const DEFAULT_PAGE = 1;
export const DEFAULT_LIMIT = 12;

/** Source type identifiers */
export const SOURCE_TYPES = {
  HACKERNEWS: 'hackernews',
  RSS: 'rss',
};

/** Human-readable source labels */
export const SOURCE_LABELS = {
  [SOURCE_TYPES.HACKERNEWS]: 'HackerNews',
  [SOURCE_TYPES.RSS]: 'RSS',
};

/** Source color mapping for badges */
export const SOURCE_COLORS = {
  [SOURCE_TYPES.HACKERNEWS]: 'var(--source-hn)',
  [SOURCE_TYPES.RSS]: 'var(--source-rss)',
};

/** Filter options for article lists */
export const FILTER_ALL = 'all';

/** Cache stale time for TanStack Query (5 minutes) */
export const QUERY_STALE_TIME = 5 * 60 * 1000;

/** Query retry count */
export const QUERY_RETRY_COUNT = 2;
