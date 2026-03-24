import { request } from './client';

export async function getFlights() {
  const data = await request('/flights');
  if (Array.isArray(data)) return data;
  return data?.flights ?? [];
}

export async function getFlight(id) {
  const data = await request(`/flights/${id}`);
  return data?.flight ?? data;
}

export function createFlight(data) {
  return request('/flights', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function updateFlight(id, data) {
  return request(`/flights/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export function deleteFlight(id) {
  return request(`/flights/${id}`, { method: 'DELETE' });
}
