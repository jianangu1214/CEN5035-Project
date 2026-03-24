import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import userEvent from '@testing-library/user-event'
import HotelForm from './HotelForm'

vi.mock('../api/hotels', () => ({
  getHotel: vi.fn(),
  createHotel: vi.fn(),
  updateHotel: vi.fn(),
}))

import { createHotel } from '../api/hotels'

describe('HotelForm', () => {
  beforeEach(() => {
    vi.mocked(createHotel).mockReset()
  })

  it('submits a new hotel via createHotel', async () => {
    const user = userEvent.setup()
    vi.mocked(createHotel).mockResolvedValue({ id: 1 })

    render(
      <MemoryRouter initialEntries={['/hotels/new']}>
        <Routes>
          <Route path="/hotels/new" element={<HotelForm />} />
          <Route path="/hotels" element={<div data-testid="after">Hotels list</div>} />
        </Routes>
      </MemoryRouter>,
    )

    await user.type(screen.getByLabelText(/hotel name/i), 'Sea View')
    await user.type(screen.getByLabelText(/^city/i), 'Tampa')
    await user.type(screen.getByLabelText(/country/i), 'USA')
    await user.type(screen.getByLabelText(/check-in/i), '2026-05-01')
    await user.type(screen.getByLabelText(/check-out/i), '2026-05-04')
    await user.clear(screen.getByLabelText(/total price/i))
    await user.type(screen.getByLabelText(/total price/i), '250')
    await user.click(screen.getByTestId('hotel-submit'))

    await waitFor(() => {
      expect(createHotel).toHaveBeenCalledWith(
        expect.objectContaining({
          hotel_name: 'Sea View',
          city: 'Tampa',
          country: 'USA',
          check_in: '2026-05-01',
          check_out: '2026-05-04',
          price: 250,
        }),
      )
    })
    await waitFor(() => {
      expect(screen.getByTestId('after')).toBeInTheDocument()
    })
  })
})
