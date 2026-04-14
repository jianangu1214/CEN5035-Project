import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import Summary from './Summary'

vi.mock('../api/summary', () => ({
  getSummary: vi.fn(),
}))

import { getSummary } from '../api/summary'

describe('Summary', () => {
  beforeEach(() => {
    vi.mocked(getSummary).mockReset()
  })

  it('loads and displays month summary by default', async () => {
    vi.mocked(getSummary).mockResolvedValue([
      {
        period: '2026-04',
        flights: 2,
        hotels: 1,
        nights: 3,
        spend: 500,
      },
    ])

    render(<Summary />)

    expect(screen.getByText('Summary')).toBeInTheDocument()

    await waitFor(() => {
      expect(getSummary).toHaveBeenCalledWith('month')
    })

    expect(await screen.findByText('2026-04')).toBeInTheDocument()
    expect(screen.getByText('2')).toBeInTheDocument()
    expect(screen.getByText('1')).toBeInTheDocument()
    expect(screen.getByText('3')).toBeInTheDocument()
    expect(screen.getByText('500')).toBeInTheDocument()
  })

  it('changes summary type to year and updates the table', async () => {
    const user = userEvent.setup()

    vi.mocked(getSummary)
      .mockResolvedValueOnce([
        {
          period: '2026-04',
          flights: 2,
          hotels: 1,
          nights: 3,
          spend: 500,
        },
      ])
      .mockResolvedValueOnce([
        {
          period: '2026',
          flights: 8,
          hotels: 4,
          nights: 10,
          spend: 2400,
        },
      ])

    render(<Summary />)

    const select = await screen.findByRole('combobox')
    await user.selectOptions(select, 'year')

    await waitFor(() => {
      expect(getSummary).toHaveBeenLastCalledWith('year')
    })

    expect(await screen.findByText('2026')).toBeInTheDocument()
    expect(screen.getByText('8')).toBeInTheDocument()
    expect(screen.getByText('4')).toBeInTheDocument()
    expect(screen.getByText('10')).toBeInTheDocument()
    expect(screen.getByText('2400')).toBeInTheDocument()
  })
})