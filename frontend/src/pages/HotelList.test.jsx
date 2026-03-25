import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import userEvent from '@testing-library/user-event'
import HotelList from './HotelList'

vi.mock('../api/hotels', () => ({
  getHotels: vi.fn(),
  deleteHotel: vi.fn(),
}))

import { getHotels, deleteHotel } from '../api/hotels'

function renderHotelList() {
  return render(
    <BrowserRouter>
      <HotelList />
    </BrowserRouter>,
  )
}

describe('HotelList', () => {
  beforeEach(() => {
    vi.mocked(getHotels).mockReset()
    vi.mocked(deleteHotel).mockReset()
    window.confirm = vi.fn(() => true)
  })

  it('shows loading then renders hotels from the API', async () => {
    vi.mocked(getHotels).mockResolvedValue([
      {
        id: 1,
        hotel_name: 'Grand',
        city: 'Miami',
        country: 'USA',
        check_in: '2024-01-01',
        check_out: '2024-01-05',
        price: 100,
        notes: 'Nice stay',
      },
    ])
    renderHotelList()
    expect(screen.getByText(/loading hotels/i)).toBeInTheDocument()
    await waitFor(() => {
      expect(screen.getByText('Grand')).toBeInTheDocument()
    })
    expect(screen.getByText('Miami')).toBeInTheDocument()
    expect(screen.getByText('$100.00')).toBeInTheDocument()
  })

  it('calls deleteHotel when delete is confirmed', async () => {
    const user = userEvent.setup()
    vi.mocked(getHotels).mockResolvedValue([
      { id: 7, hotel_name: 'ToRemove', city: 'X', country: 'Y', check_in: '2024-01-01', check_out: '2024-01-02', price: 1, notes: '' },
    ])
    vi.mocked(deleteHotel).mockResolvedValue(undefined)
    renderHotelList()
    await waitFor(() => expect(screen.getByText('ToRemove')).toBeInTheDocument())
    await user.click(screen.getByRole('button', { name: /delete/i }))
    expect(deleteHotel).toHaveBeenCalledWith(7)
    await waitFor(() => {
      expect(screen.queryByText('ToRemove')).not.toBeInTheDocument()
    })
  })
})
