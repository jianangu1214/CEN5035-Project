import { useState, useEffect, useMemo } from 'react'
import { MapContainer, TileLayer, Marker, Popup, Polyline, useMap } from 'react-leaflet'
import L from 'leaflet'
import { getHotels } from '../api/hotels'
import { getFlights } from '../api/flights'
import { getAirportCoords } from '../data/airports'
import { geocodeCity } from '../utils/geocode'

import 'leaflet/dist/leaflet.css'

const hotelIcon = new L.Icon({
  iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-blue.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41],
})

const airportIcon = new L.Icon({
  iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41],
})

function FitBounds({ positions }) {
  const map = useMap()
  useEffect(() => {
    if (positions.length === 0) return
    const bounds = L.latLngBounds(positions)
    map.fitBounds(bounds, { padding: [50, 50], maxZoom: 12 })
  }, [map, positions])
  return null
}

function formatDate(str) {
  if (!str) return ''
  const d = new Date(str)
  return isNaN(d.getTime()) ? str : d.toLocaleDateString()
}

function formatDateTime(str) {
  if (!str) return ''
  const d = new Date(str)
  return isNaN(d.getTime()) ? str : d.toLocaleString()
}

function formatPrice(n) {
  if (n == null || n === '') return ''
  const num = Number(n)
  return isNaN(num) ? n : `$${num.toFixed(2)}`
}

export default function TravelMap() {
  const [hotels, setHotels] = useState([])
  const [flights, setFlights] = useState([])
  const [hotelCoords, setHotelCoords] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [showHotels, setShowHotels] = useState(true)
  const [showFlights, setShowFlights] = useState(true)

  useEffect(() => {
    let cancelled = false
    async function load() {
      setLoading(true)
      setError(null)
      try {
        const [hotelData, flightData] = await Promise.all([getHotels(), getFlights()])
        if (cancelled) return

        const h = Array.isArray(hotelData) ? hotelData : []
        const f = Array.isArray(flightData) ? flightData : []
        setHotels(h)
        setFlights(f)

        const coords = await Promise.all(
          h.map(async (hotel) => {
            const c = await geocodeCity(hotel.city, hotel.country)
            return c ? { ...hotel, coords: c } : null
          })
        )
        if (!cancelled) {
          setHotelCoords(coords.filter(Boolean))
        }
      } catch (err) {
        if (!cancelled) {
          setError(err.body?.error || err.message || 'Failed to load data')
        }
      } finally {
        if (!cancelled) setLoading(false)
      }
    }
    load()
    return () => { cancelled = true }
  }, [])

  const flightRoutes = useMemo(() => {
    return flights
      .map((f) => {
        const from = getAirportCoords(f.from_airport)
        const to = getAirportCoords(f.to_airport)
        if (!from || !to) return null
        return { ...f, from, to }
      })
      .filter(Boolean)
  }, [flights])

  const allPositions = useMemo(() => {
    const pts = []
    if (showHotels) {
      hotelCoords.forEach((h) => pts.push([h.coords.lat, h.coords.lng]))
    }
    if (showFlights) {
      flightRoutes.forEach((r) => {
        pts.push([r.from.lat, r.from.lng])
        pts.push([r.to.lat, r.to.lng])
      })
    }
    return pts
  }, [hotelCoords, flightRoutes, showHotels, showFlights])

  const unmatchedFlights = useMemo(
    () =>
      flights.filter((f) => {
        const from = getAirportCoords(f.from_airport)
        const to = getAirportCoords(f.to_airport)
        return !from || !to
      }),
    [flights]
  )

  if (loading) {
    return <div className="loading">Loading map data…</div>
  }

  return (
    <>
      <div className="page-header">
        <h1>Travel Map</h1>
        <div className="map-toggles">
          <label className="toggle-label">
            <input
              type="checkbox"
              checked={showHotels}
              onChange={(e) => setShowHotels(e.target.checked)}
            />
            <span className="toggle-dot toggle-dot--hotel" />
            Hotels ({hotelCoords.length})
          </label>
          <label className="toggle-label">
            <input
              type="checkbox"
              checked={showFlights}
              onChange={(e) => setShowFlights(e.target.checked)}
            />
            <span className="toggle-dot toggle-dot--flight" />
            Flights ({flightRoutes.length})
          </label>
        </div>
      </div>

      {error && (
        <div className="error-banner" role="alert">
          {error}
        </div>
      )}

      {unmatchedFlights.length > 0 && (
        <div className="map-info-banner">
          {unmatchedFlights.length} flight{unmatchedFlights.length > 1 ? 's' : ''} could not be
          mapped (unrecognized airport code
          {unmatchedFlights.length > 1 ? 's' : ''}:{' '}
          {unmatchedFlights
            .flatMap((f) => {
              const codes = []
              if (!getAirportCoords(f.from_airport)) codes.push(f.from_airport)
              if (!getAirportCoords(f.to_airport)) codes.push(f.to_airport)
              return codes
            })
            .filter((v, i, a) => a.indexOf(v) === i)
            .join(', ')}
          )
        </div>
      )}

      <div className="map-container">
        <MapContainer
          center={[20, 0]}
          zoom={2}
          scrollWheelZoom
          className="travel-map"
        >
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />

          {allPositions.length > 0 && <FitBounds positions={allPositions} />}

          {showHotels &&
            hotelCoords.map((h) => (
              <Marker
                key={`hotel-${h.id}`}
                position={[h.coords.lat, h.coords.lng]}
                icon={hotelIcon}
              >
                <Popup>
                  <div className="map-popup">
                    <strong>{h.hotel_name || 'Hotel'}</strong>
                    <br />
                    {h.city}{h.country ? `, ${h.country}` : ''}
                    {h.check_in && (
                      <>
                        <br />
                        {formatDate(h.check_in)} – {formatDate(h.check_out)}
                      </>
                    )}
                    {h.price != null && h.price !== '' && (
                      <>
                        <br />
                        {formatPrice(h.price)}
                      </>
                    )}
                  </div>
                </Popup>
              </Marker>
            ))}

          {showFlights &&
            flightRoutes.map((r) => (
              <Polyline
                key={`route-${r.id}`}
                positions={[
                  [r.from.lat, r.from.lng],
                  [r.to.lat, r.to.lng],
                ]}
                pathOptions={{
                  color: '#ef4444',
                  weight: 2.5,
                  opacity: 0.8,
                  dashArray: '8 6',
                }}
              >
                <Popup>
                  <div className="map-popup">
                    <strong>
                      {r.airline || ''} {r.flight_number || ''}
                    </strong>
                    <br />
                    {r.from_airport} → {r.to_airport}
                    {r.depart_time && (
                      <>
                        <br />
                        {formatDateTime(r.depart_time)}
                      </>
                    )}
                    {r.price != null && r.price !== '' && (
                      <>
                        <br />
                        {formatPrice(r.price)}
                      </>
                    )}
                  </div>
                </Popup>
              </Polyline>
            ))}

          {showFlights &&
            flightRoutes.map((r) => (
              <Marker
                key={`airport-from-${r.id}`}
                position={[r.from.lat, r.from.lng]}
                icon={airportIcon}
              >
                <Popup>
                  <div className="map-popup">
                    <strong>{r.from_airport}</strong>
                    <br />
                    {r.from.name}
                  </div>
                </Popup>
              </Marker>
            ))}

          {showFlights &&
            flightRoutes.map((r) => (
              <Marker
                key={`airport-to-${r.id}`}
                position={[r.to.lat, r.to.lng]}
                icon={airportIcon}
              >
                <Popup>
                  <div className="map-popup">
                    <strong>{r.to_airport}</strong>
                    <br />
                    {r.to.name}
                  </div>
                </Popup>
              </Marker>
            ))}
        </MapContainer>
      </div>

      <div className="map-legend">
        <div className="legend-item">
          <span className="legend-marker legend-marker--hotel" />
          Hotel stays
        </div>
        <div className="legend-item">
          <span className="legend-marker legend-marker--flight" />
          Airports
        </div>
        <div className="legend-item">
          <span className="legend-line" />
          Flight routes
        </div>
      </div>
    </>
  )
}
