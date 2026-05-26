import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

vi.mock('../utils/constants', () => ({
  API_BASE_URL: 'http://test-api',
  API_PREFIX: '',
  DEFAULT_LIMIT: 12,
}));

const {
  fetchArticles,
  fetchArticle,
  fetchTopics,
  fetchTopicArticles,
  searchArticles,
  generateSummary,
  createTopic,
  updateTopic,
  deleteTopic,
  triggerFetch,
  triggerProcess,
  api,
} = await import('./api');

const mockSuccessResponse = (data, status = 200) =>
  Promise.resolve({
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(data),
    text: () => Promise.resolve(JSON.stringify(data)),
  });

describe('API service', () => {
  beforeEach(() => {
    localStorage.clear();
    globalThis.fetch = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('fetchArticles', () => {
    it('fetches articles with default pagination', async () => {
      const articles = [{ id: '1', title: 'Test' }];
      globalThis.fetch.mockResolvedValue(mockSuccessResponse(articles));

      const result = await fetchArticles();

      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://test-api/articles?page=1&limit=12'
      );
      expect(result).toEqual(articles);
    });

    it('fetches articles with custom page and limit', async () => {
      globalThis.fetch.mockResolvedValue(mockSuccessResponse([]));

      await fetchArticles(3, 20);

      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://test-api/articles?page=3&limit=20'
      );
    });

    it('throws on non-ok response', async () => {
      globalThis.fetch.mockResolvedValue(mockSuccessResponse(null, 500));

      await expect(fetchArticles()).rejects.toThrow('Failed to fetch articles');
    });
  });

  describe('fetchArticle', () => {
    it('fetches a single article by ID', async () => {
      const article = { id: '42', title: 'Test' };
      globalThis.fetch.mockResolvedValue(mockSuccessResponse(article));

      const result = await fetchArticle('42');

      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://test-api/articles/42'
      );
      expect(result).toEqual(article);
    });
  });

  describe('fetchTopics', () => {
    it('fetches all topics', async () => {
      const topics = [{ id: '1', name: 'Tech' }];
      globalThis.fetch.mockResolvedValue(mockSuccessResponse(topics));

      const result = await fetchTopics();

      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://test-api/topics'
      );
      expect(result).toEqual(topics);
    });
  });

  describe('fetchTopicArticles', () => {
    it('fetches articles for a topic with pagination', async () => {
      globalThis.fetch.mockResolvedValue(mockSuccessResponse([]));

      await fetchTopicArticles('ai', 2, 10);

      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://test-api/topics/ai/articles?page=2&limit=10'
      );
    });
  });

  describe('searchArticles', () => {
    it('sends search query', async () => {
      globalThis.fetch.mockResolvedValue(mockSuccessResponse([]));

      await searchArticles('gemini');

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('q=gemini')
      );
    });

    it('includes source filter when provided', async () => {
      globalThis.fetch.mockResolvedValue(mockSuccessResponse([]));

      await searchArticles('test', 'hackernews');

      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('source=hackernews')
      );
    });

    it('omits source filter when source is Tümü', async () => {
      globalThis.fetch.mockResolvedValue(mockSuccessResponse([]));

      await searchArticles('test', 'Tümü');

      const callUrl = globalThis.fetch.mock.calls[0][0];
      expect(callUrl).not.toContain('source=');
    });
  });

  describe('generateSummary', () => {
    it('posts to generate summary', async () => {
      const summary = { summary_tr: 'Özet' };
      globalThis.fetch.mockResolvedValue(mockSuccessResponse(summary));

      const result = await generateSummary('42');

      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://test-api/articles/42/summary',
        expect.objectContaining({ method: 'POST' })
      );
      expect(result).toEqual(summary);
    });
  });

  describe('admin endpoints', () => {
    const adminKey = 'admin-secret-123';

    beforeEach(() => {
      localStorage.setItem('admin_api_key', adminKey);
    });

    it('createTopic sends POST with admin headers', async () => {
      const topic = { name: 'AI', slug: 'ai' };
      globalThis.fetch.mockResolvedValue(mockSuccessResponse(topic, 201));

      const result = await createTopic(topic);

      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://test-api/admin/topics',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'X-Admin-API-Key': adminKey,
          }),
          body: JSON.stringify(topic),
        })
      );
      expect(result).toEqual(topic);
    });

    it('updateTopic sends PUT with admin headers', async () => {
      globalThis.fetch.mockResolvedValue(mockSuccessResponse({}));

      await updateTopic('1', { name: 'Updated' });

      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://test-api/admin/topics/1',
        expect.objectContaining({ method: 'PUT' })
      );
    });

    it('deleteTopic sends DELETE', async () => {
      globalThis.fetch.mockResolvedValue(mockSuccessResponse({}));

      await deleteTopic('1');

      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://test-api/admin/topics/1',
        expect.objectContaining({ method: 'DELETE' })
      );
    });

    it('triggerFetch sends POST to admin/fetch', async () => {
      globalThis.fetch.mockResolvedValue(mockSuccessResponse({ fetched: 5 }));

      const result = await triggerFetch();

      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://test-api/admin/fetch',
        expect.objectContaining({ method: 'POST' })
      );
      expect(result).toEqual({ fetched: 5 });
    });

    it('triggerProcess sends POST to admin/process', async () => {
      globalThis.fetch.mockResolvedValue(mockSuccessResponse({ processed: 3 }));

      const result = await triggerProcess();

      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://test-api/admin/process',
        expect.objectContaining({ method: 'POST' })
      );
      expect(result).toEqual({ processed: 3 });
    });

    it('redirects to login on 401', async () => {
      delete window.location;
      window.location = { href: '' };
      globalThis.fetch.mockResolvedValue(mockSuccessResponse(null, 401));

      await expect(createTopic({ name: 'test' })).rejects.toThrow('Unauthorized');
      expect(localStorage.getItem('admin_api_key')).toBeNull();
      expect(window.location.href).toBe('/login');
    });
  });

  describe('api convenience object', () => {
    it('exports all functions', () => {
      expect(api.getArticles).toBe(fetchArticles);
      expect(api.getArticle).toBe(fetchArticle);
      expect(api.getTopics).toBe(fetchTopics);
      expect(api.getTopicArticles).toBe(fetchTopicArticles);
      expect(api.searchArticles).toBe(searchArticles);
      expect(api.generateSummary).toBe(generateSummary);
      expect(api.createTopic).toBe(createTopic);
      expect(api.updateTopic).toBe(updateTopic);
      expect(api.deleteTopic).toBe(deleteTopic);
      expect(api.triggerFetch).toBe(triggerFetch);
      expect(api.triggerProcess).toBe(triggerProcess);
    });
  });
});
