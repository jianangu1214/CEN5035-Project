import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import TravelMap from './TravelMap'

vi.mock('../api/hotels', () => ({
  getHotels: vi.fn(),
}))

vi.mock('../api/flights', () => ({
  getFlights: vi.fn(),
}))

vi.mock('../utils/geocode', () => ({
  geocodeCity: vi.fn(),
}))

vi.mock('../data/airports', () => ({
  getAirportCoords: vi.fn(),
}))

vi.mock('react-leaflet', () => ({
  MapContainer: ({ children }) => <div data-testid="map-container">{children}</div>,
  TileLayer: () => <div data-testid="tile-layer" />,
  Marker: ({ children }) => <div data-testid="marker">{children}</div>,
  Popup: ({ children }) => <div data-testid="popup">{children}</div>,
  Polyline: ({ children }) => <div data-testid="polyline">{children}</div>,
  useMap: () => ({
    fitBounds: vi.fn(),
  }),
}))

vi.mock('leaflet', () => {
  class Icon {
    constructor(options) {
      this.options = options
    }
  }

  return {
    default: {
      Icon,
      latLngBounds: vi.fn(() => ({})),
    },
    Icon,
    latLngBounds: vi.fn(() => ({})),
  }
})

import { getHotels } from '../api/hotels'
import { getFlights } from '../api/flights'
import { geocodeCity } from '../utils/geocode'
import { getAirportCoords } from '../data/airports'

describe('TravelMap', () => {
  beforeEach(() => {
    vi.mocked(getHotels).mockReset()
    vi.mocked(getFlights).mockReset()
    vi.mocked(geocodeCity).mockReset()
    vi.mocked(getAirportCoords).mockReset()
  })

  it('loads and displays mapped hotel and flight counts', async () => {
    vi.mocked(getHotels).mockResolvedValue([
      {
        id: 1,
        hotel_name: 'Hilton Tokyo',
        city: 'Tokyo',
        country: 'Japan',
        check_in: '2026-04-01',
        check_out: '2026-04-03',
        price: 300,
      },
    ])

    vi.mocked(getFlights).mockResolvedValue([
      {
        id: 10,
        airline: 'ANA',
        flight_number: 'NH12',
        from_airport: 'JFK',
        to_airport: 'LAX',
        depart_time: '2026-04-01T10:00:00Z',
        arrive_time: '2026-04-01T14:00:00Z',
        price: 200,
      },
    ])

    vi.mocked(geocodeCity).mockResolvedValue({
      lat: 35.6762,
      lng: 139.6503,
    })

    vi.mocked(getAirportCoords).mockImplementation((code) => {
      if (code === 'JFK') {
        return {
          lat: 40.6413,
          lng: -73.7781,
          name: 'John F. Kennedy International Airport',
        }
      }
      if (code === 'LAX') {
        return {
          lat: 33.9416,
          lng: -118.4085,
          name: 'Los Angeles International Airport',
        }
      }
      return null
    })

    render(<TravelMap />)

    expect(screen.getByText('Loading map data…')).toBeInTheDocument()

    expect(await screen.findByText('Travel Map')).toBeInTheDocument()
    expect(screen.getByText(/Hotels \(1\)/)).toBeInTheDocument()
    expect(screen.getByText(/Flights \(1\)/)).toBeInTheDocument()
  })

  it('shows unmatched flight warning when airport codes are unknown', async () => {
    vi.mocked(getHotels).mockResolvedValue([])

    vi.mocked(getFlights).mockResolvedValue([
      {
        id: 20,
        airline: 'Test Air',
        flight_number: 'TA100',
        from_airport: 'XXX',
        to_airport: 'YYY',
        depart_time: '2026-04-01T10:00:00Z',
        arrive_time: '2026-04-01T14:00:00Z',
        price: 120,
      },
    ])

    vi.mocked(geocodeCity).mockResolvedValue(null)
    vi.mocked(getAirportCoords).mockReturnValue(null)

    render(<TravelMap />)

    expect(await screen.findByText('Travel Map')).toBeInTheDocument()

    await waitFor(() => {
      expect(screen.getByText(/could not be mapped/i)).toBeInTheDocument()
    })
  })
})