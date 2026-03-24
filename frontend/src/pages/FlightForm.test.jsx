import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import userEvent from '@testing-library/user-event'
import FlightForm from './FlightForm'

vi.mock('../api/flights', () => ({
  getFlight: vi.fn(),
  createFlight: vi.fn(),
  updateFlight: vi.fn(),
}))

import { createFlight } from '../api/flights'

describe('FlightForm', () => {
  beforeEach(() => {
    vi.mocked(createFlight).mockReset()
  })

  it('creates a flight with RFC3339 times', async () => {
    const user = userEvent.setup()
    vi.mocked(createFlight).mockResolvedValue({})

    render(
      <MemoryRouter initialEntries={['/flights/new']}>
        <Routes>
          <Route path="/flights/new" element={<FlightForm />} />
          <Route path="/flights" element={<div data-testid="flight-list">Flights</div>} />
        </Routes>
      </MemoryRouter>,
    )

    await user.type(screen.getByTestId('flight-airline'), 'UnitAir')
    await user.type(screen.getByTestId('flight-number'), 'UA100')
    await user.type(screen.getByTestId('flight-from'), 'mia')
    await user.type(screen.getByTestId('flight-to'), 'jfk')
    await user.type(screen.getByTestId('flight-depart'), '2026-07-10T09:00')
    await user.type(screen.getByTestId('flight-arrive'), '2026-07-10T12:30')
    await user.clear(screen.getByTestId('flight-price'))
    await user.type(screen.getByTestId('flight-price'), '320')
    await user.click(screen.getByTestId('flight-submit'))

    await waitFor(() => {
      expect(createFlight).toHaveBeenCalled()
    })
    const body = vi.mocked(createFlight).mock.calls[0][0]
    expect(body).toMatchObject({
      airline: 'UnitAir',
      flight_number: 'UA100',
      from_airport: 'MIA',
      to_airport: 'JFK',
      price: 320,
    })
    expect(body.depart_time).toMatch(/^\d{4}-\d{2}-\d{2}T/)
    expect(body.arrive_time).toMatch(/^\d{4}-\d{2}-\d{2}T/)

    await waitFor(() => {
      expect(screen.getByTestId('flight-list')).toBeInTheDocument()
    })
  })
})
