const cache = new Map();

export async function geocodeCity(city, country) {
  if (!city) return null;

  const key = `${city}|${country || ''}`.toLowerCase();
  if (cache.has(key)) return cache.get(key);

  const q = country ? `${city}, ${country}` : city;
  try {
    const res = await fetch(
      `https://nominatim.openstreetmap.org/search?` +
        new URLSearchParams({ q, format: 'json', limit: '1' }),
      { headers: { 'User-Agent': 'TravelLog/1.0' } }
    );
    if (!res.ok) return null;
    const data = await res.json();
    if (!data.length) return null;

    const coords = { lat: parseFloat(data[0].lat), lng: parseFloat(data[0].lon) };
    cache.set(key, coords);
    return coords;
  } catch {
    return null;
  }
}
