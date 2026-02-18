/**
 * API client for TravelLog backend.
 * Backend runs on port 8080; hotel endpoints require JWT in Authorization header.
 * Store token after login: localStorage.setItem('travellog_token', token)
 */
const API_BASE = import.meta.env.VITE_API_URL || '/api';

const AUTH_TOKEN_KEY = 'travellog_token';

export function getAuthToken() {
  return localStorage.getItem(AUTH_TOKEN_KEY);
}

export function setAuthToken(token) {
  if (token) localStorage.setItem(AUTH_TOKEN_KEY, token);
  else localStorage.removeItem(AUTH_TOKEN_KEY);
}

export async function request(path, options = {}) {
  const url = `${API_BASE}${path}`;
  const token = getAuthToken();
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const res = await fetch(url, { ...options, headers });

  if (!res.ok) {
    const err = new Error(res.status === 401 ? 'Please log in' : (res.statusText || 'Request failed'));
    err.status = res.status;
    try {
      err.body = await res.json();
      if (err.body.error || err.body.message) {
        err.message = err.body.error || err.body.message;
      }
    } catch {
      try {
        err.body = await res.text();
      } catch {
        err.body = null;
      }
    }
    throw err;
  }

  const contentType = res.headers.get('content-type');
  if (contentType && contentType.includes('application/json')) {
    return res.json();
  }
  return res.text();
}
