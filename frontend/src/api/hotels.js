import { request } from './client';

export function getHotels() {
  return request('/hotels');
}

export function getHotel(id) {
  return request(`/hotels/${id}`);
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
