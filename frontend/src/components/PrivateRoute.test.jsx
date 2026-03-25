import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import PrivateRoute from './PrivateRoute'

vi.mock('../api/client', () => ({
  getAuthToken: vi.fn(),
}))

import { getAuthToken } from '../api/client'

describe('PrivateRoute', () => {
  beforeEach(() => {
    vi.mocked(getAuthToken).mockReset()
  })

  it('redirects to /login when there is no token', () => {
    vi.mocked(getAuthToken).mockReturnValue(null)
    render(
      <MemoryRouter initialEntries={['/hotels']}>
        <Routes>
          <Route element={<PrivateRoute />}>
            <Route path="/hotels" element={<div data-testid="protected">Secret</div>} />
          </Route>
          <Route path="/login" element={<div data-testid="login-page">Login page</div>} />
        </Routes>
      </MemoryRouter>,
    )
    expect(screen.getByTestId('login-page')).toBeInTheDocument()
    expect(screen.queryByTestId('protected')).not.toBeInTheDocument()
  })

  it('renders nested route when a token is present', () => {
    vi.mocked(getAuthToken).mockReturnValue('test-jwt-token')
    render(
      <MemoryRouter initialEntries={['/hotels']}>
        <Routes>
          <Route element={<PrivateRoute />}>
            <Route path="/hotels" element={<div data-testid="protected">Secret</div>} />
          </Route>
          <Route path="/login" element={<div data-testid="login-page">Login page</div>} />
        </Routes>
      </MemoryRouter>,
    )
    expect(screen.getByTestId('protected')).toHaveTextContent('Secret')
    expect(screen.queryByTestId('login-page')).not.toBeInTheDocument()
  })
})
