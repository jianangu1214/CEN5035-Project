import { request } from './client';

export async function getHotels() {
  const data = await request('/hotels');
  if (Array.isArray(data)) return data;
  return data?.hotels ?? [];
}

export async function getHotel(id) {
  const data = await request(`/hotels/${id}`);
  return data?.hotel ?? data;
}

export function createHotel(data) {
  return request('/hotels', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function updateHotel(id, data) {
  return request(`/hotels/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export function deleteHotel(id) {
  return request(`/hotels/${id}`, { method: 'DELETE' });
}
