import { API_BASE_URL as CONST_BASE_URL, API_PREFIX, DEFAULT_LIMIT } from '../utils/constants';

const API_BASE_URL = `${CONST_BASE_URL}${API_PREFIX}`;

function getHeaders(isAdmin = false) {
  const headers = {
    'Content-Type': 'application/json',
  };
  if (isAdmin) {
    const adminKey = localStorage.getItem('admin_api_key');
    if (adminKey) {
      headers['X-Admin-API-Key'] = adminKey;
    }
  }
  return headers;
}

export async function fetchArticles(page = 1, limit = DEFAULT_LIMIT) {
  const response = await fetch(`${API_BASE_URL}/articles?page=${page}&limit=${limit}`);
  if (!response.ok) throw new Error('Failed to fetch articles');
  return response.json();
}

export async function fetchArticle(id) {
  const response = await fetch(`${API_BASE_URL}/articles/${id}`);
  if (!response.ok) throw new Error('Failed to fetch article');
  return response.json();
}

export async function fetchTopics() {
  const response = await fetch(`${API_BASE_URL}/topics`);
  if (!response.ok) throw new Error('Failed to fetch topics');
  return response.json();
}

export async function fetchTopicArticles(slug, page = 1, limit = DEFAULT_LIMIT) {
  const response = await fetch(`${API_BASE_URL}/topics/${slug}/articles?page=${page}&limit=${limit}`);
  if (!response.ok) throw new Error('Failed to fetch topic articles');
  return response.json();
}

export async function searchArticles(query, source = '', page = 1, limit = DEFAULT_LIMIT) {
  const params = new URLSearchParams();
  if (query) params.append('q', query);
  if (source && source !== 'Tümü') params.append('source', source);
  params.append('page', String(page));
  params.append('limit', String(limit));
  
  const response = await fetch(`${API_BASE_URL}/articles/search?${params.toString()}`);
  if (!response.ok) throw new Error('Arama yapılamadı');
  return response.json();
}

export async function generateSummary(id) {
  const response = await fetch(`${API_BASE_URL}/articles/${id}/summary`, {
    method: 'POST',
    headers: getHeaders(false),
  });
  if (!response.ok) throw new Error('Failed to generate summary');
  return response.json();
}

// Admin API calls
export async function createTopic(topic) {
  const response = await fetch(`${API_BASE_URL}/admin/topics`, {
    method: 'POST',
    headers: getHeaders(true),
    body: JSON.stringify(topic),
  });
  if (response.status === 401) {
    localStorage.removeItem('admin_api_key');
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }
  if (!response.ok) throw new Error('Failed to create topic');
  return response.json();
}

export async function updateTopic(id, topic) {
  const response = await fetch(`${API_BASE_URL}/admin/topics/${id}`, {
    method: 'PUT',
    headers: getHeaders(true),
    body: JSON.stringify(topic),
  });
  if (response.status === 401) {
    localStorage.removeItem('admin_api_key');
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }
  if (!response.ok) throw new Error('Failed to update topic');
  return response.json();
}

export async function deleteTopic(id) {
  const response = await fetch(`${API_BASE_URL}/admin/topics/${id}`, {
    method: 'DELETE',
    headers: getHeaders(true),
  });
  if (response.status === 401) {
    localStorage.removeItem('admin_api_key');
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }
  if (!response.ok) throw new Error('Failed to delete topic');
  return response.text();
}

export async function triggerFetch() {
  const response = await fetch(`${API_BASE_URL}/admin/fetch`, {
    method: 'POST',
    headers: getHeaders(true),
  });
  if (response.status === 401) {
    localStorage.removeItem('admin_api_key');
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }
  if (!response.ok) throw new Error('Failed to trigger fetch');
  return response.json();
}

export async function triggerProcess() {
  const response = await fetch(`${API_BASE_URL}/admin/process`, {
    method: 'POST',
    headers: getHeaders(true),
  });
  if (response.status === 401) {
    localStorage.removeItem('admin_api_key');
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }
  if (!response.ok) throw new Error('Failed to trigger process');
  return response.json();
}

export const api = {
  getArticles: fetchArticles,
  getArticle: fetchArticle,
  getTopics: fetchTopics,
  getTopicArticles: fetchTopicArticles,
  searchArticles: searchArticles,
  generateSummary: generateSummary,
  createTopic: createTopic,
  updateTopic: updateTopic,
  deleteTopic: deleteTopic,
  triggerFetch: triggerFetch,
  triggerProcess: triggerProcess,
};
