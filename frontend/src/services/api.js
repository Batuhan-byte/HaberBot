import { API_BASE_URL as CONST_BASE_URL, API_PREFIX, DEFAULT_LIMIT } from '../utils/constants';

const API_BASE_URL = `${CONST_BASE_URL}${API_PREFIX}`;

function getHeaders(isAdmin = false, isMultipart = false) {
  const headers = {};
  if (!isMultipart) {
    headers['Content-Type'] = 'application/json';
  }
  const token = localStorage.getItem('access_token');
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
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
  const response = await fetch(`${API_BASE_URL}/articles/${id}`, {
    headers: getHeaders(true),
  });
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

export async function getArticleSummary(id) {
  const response = await fetch(`${API_BASE_URL}/articles/${id}/summary`, {
    method: 'GET',
    headers: getHeaders(false),
  });
  if (!response.ok) throw new Error('Failed to fetch summary');
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

export async function fetchArticlesAdmin(topicId = '', page = 1, limit = 50) {
  const params = new URLSearchParams();
  if (topicId) params.append('topic_id', topicId);
  params.append('page', String(page));
  params.append('limit', String(limit));

  const response = await fetch(`${API_BASE_URL}/admin/articles?${params.toString()}`, {
    method: 'GET',
    headers: getHeaders(true),
  });
  if (response.status === 401) {
    localStorage.removeItem('admin_api_key');
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }
  if (!response.ok) throw new Error('Failed to fetch admin articles');
  return response.json();
}

export async function updateArticleAdmin(id, updates) {
  const response = await fetch(`${API_BASE_URL}/admin/articles/${id}`, {
    method: 'PUT',
    headers: getHeaders(true),
    body: JSON.stringify(updates),
  });
  if (response.status === 401) {
    localStorage.removeItem('admin_api_key');
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }
  if (!response.ok) throw new Error('Failed to update article');
  return response.json();
}

export async function deleteArticleAdmin(id) {
  const response = await fetch(`${API_BASE_URL}/admin/articles/${id}`, {
    method: 'DELETE',
    headers: getHeaders(true),
  });
  if (response.status === 401) {
    localStorage.removeItem('admin_api_key');
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }
  if (!response.ok) throw new Error('Failed to delete article');
  return response.json();
}

// User Authentication API Calls
export async function login(username, password) {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.error || 'Giriş başarısız.');
  }
  const data = await response.json();
  localStorage.setItem('access_token', data.access_token);
  localStorage.setItem('user', JSON.stringify(data.user));
  if (data.user.role === 'Admin') {
    localStorage.setItem('admin_api_key', 'haberbot-super-secret-admin-key');
  }
  return data;
}

export async function register(username, email, password) {
  const response = await fetch(`${API_BASE_URL}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, email, password }),
  });
  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.error || 'Kayıt başarısız.');
  }
  return response.json();
}

export async function logout() {
  try {
    await fetch(`${API_BASE_URL}/auth/logout`, {
      method: 'POST',
      headers: getHeaders(false),
    });
  } catch (e) {
    console.error("Logout request failed:", e);
  }
  localStorage.removeItem('access_token');
  localStorage.removeItem('user');
  localStorage.removeItem('admin_api_key');
}

// Comments API Calls
export async function fetchComments(articleId) {
  const response = await fetch(`${API_BASE_URL}/comments?article_id=${articleId}`);
  if (!response.ok) throw new Error('Yorumlar yüklenemedi.');
  return response.json();
}

export async function createComment(articleId, content) {
  const response = await fetch(`${API_BASE_URL}/comments`, {
    method: 'POST',
    headers: getHeaders(false),
    body: JSON.stringify({ article_id: articleId, content }),
  });
  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.error || 'Yorum gönderilemedi.');
  }
  return response.json();
}

// Profile API Calls
export async function fetchUserProfile(username) {
  const response = await fetch(`${API_BASE_URL}/users/${username}`, {
    headers: getHeaders(false),
  });
  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.error || 'Profil yüklenemedi.');
  }
  return response.json();
}

export async function updateUserProfile(formData) {
  const response = await fetch(`${API_BASE_URL}/users/profile`, {
    method: 'PUT',
    headers: getHeaders(false, true), // isMultipart = true
    body: formData,
  });
  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.error || 'Profil güncellenemedi.');
  }
  return response.json();
}

export async function addFavorite(favoriteUserId) {
  const response = await fetch(`${API_BASE_URL}/users/${favoriteUserId}/favorites`, {
    method: 'POST',
    headers: getHeaders(false),
  });
  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.error || 'Favoriye eklenemedi.');
  }
  return response.json();
}

export async function removeFavorite(favoriteUserId) {
  const response = await fetch(`${API_BASE_URL}/users/${favoriteUserId}/favorites`, {
    method: 'DELETE',
    headers: getHeaders(false),
  });
  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.error || 'Favorilerden çıkarılamadı.');
  }
  return response;
}

export async function fetchFavorites(page = 1, limit = 20) {
  const response = await fetch(`${API_BASE_URL}/users/me/favorites?page=${page}&limit=${limit}`, {
    headers: getHeaders(false),
  });
  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.error || 'Favoriler yüklenemedi.');
  }
  return response.json();
}

export async function reportProfile(reportedUserId, reason, comment) {
  const response = await fetch(`${API_BASE_URL}/profile-reports`, {
    method: 'POST',
    headers: getHeaders(false),
    body: JSON.stringify({ reported_user_id: reportedUserId, reason, comment }),
  });
  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.error || 'Şikayet iletilemedi.');
  }
  return response.json();
}

export const api = {
  getArticles: fetchArticles,
  getArticle: fetchArticle,
  getTopics: fetchTopics,
  getTopicArticles: fetchTopicArticles,
  searchArticles: searchArticles,
  generateSummary: generateSummary,
  getArticleSummary: getArticleSummary,
  createTopic: createTopic,
  updateTopic: updateTopic,
  deleteTopic: deleteTopic,
  triggerFetch: triggerFetch,
  triggerProcess: triggerProcess,
  getArticlesAdmin: fetchArticlesAdmin,
  updateArticleAdmin: updateArticleAdmin,
  deleteArticleAdmin: deleteArticleAdmin,
  login,
  register,
  logout,
  getComments: fetchComments,
  createComment,
  fetchUserProfile,
  updateUserProfile,
  addFavorite,
  removeFavorite,
  fetchFavorites,
  reportProfile,
};

