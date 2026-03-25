import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import userEvent from '@testing-library/user-event'
import FlightList from './FlightList'

vi.mock('../api/flights', () => ({
  getFlights: vi.fn(),
  deleteFlight: vi.fn(),
}))

import { getFlights, deleteFlight } from '../api/flights'

function renderFlightList() {
  return render(
    <BrowserRouter>
      <FlightList />
    </BrowserRouter>,
  )
}

describe('FlightList', () => {
  beforeEach(() => {
    vi.mocked(getFlights).mockReset()
    vi.mocked(deleteFlight).mockReset()
    window.confirm = vi.fn(() => true)
  })

  it('loads and displays flights', async () => {
    vi.mocked(getFlights).mockResolvedValue([
      {
        id: 2,
        airline: 'TestAir',
        flight_number: 'TA99',
        from_airport: 'MIA',
        to_airport: 'JFK',
        depart_time: '2026-06-01T10:00:00Z',
        arrive_time: '2026-06-01T14:00:00Z',
        price: 199,
        notes: '',
      },
    ])
    renderFlightList()
    expect(screen.getByText(/loading flights/i)).toBeInTheDocument()
    await waitFor(() => {
      expect(screen.getByText('TestAir')).toBeInTheDocument()
    })
    expect(screen.getByText('TA99')).toBeInTheDocument()
    expect(screen.getByText('MIA')).toBeInTheDocument()
  })

  it('deletes a flight when confirmed', async () => {
    const user = userEvent.setup()
    vi.mocked(getFlights).mockResolvedValue([
      {
        id: 3,
        airline: 'X',
        flight_number: '1',
        from_airport: 'AAA',
        to_airport: 'BBB',
        depart_time: '2026-01-01T00:00:00Z',
        arrive_time: '2026-01-01T01:00:00Z',
        price: 1,
        notes: '',
      },
    ])
    vi.mocked(deleteFlight).mockResolvedValue(undefined)
    renderFlightList()
    await waitFor(() => expect(screen.getByText('X')).toBeInTheDocument())
    await user.click(screen.getByRole('button', { name: /delete/i }))
    expect(deleteFlight).toHaveBeenCalledWith(3)
  })
})
