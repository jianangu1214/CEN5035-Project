import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { request, setAuthToken } from '../api/client'

export default function Login() {
  const nav = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')

  async function onSubmit(e) {
    e.preventDefault()
    setError('')
    try {
      const data = await request('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      })

      // 兼容不同字段名
      const token = data?.token || data?.access_token || data?.jwt
      if (!token) throw new Error('No token returned from server')

      setAuthToken(token)
      nav('/hotels')
    } catch (err) {
      setError(err?.message || 'Login failed')
    }
  }

  return (
    <div>
      <h2>Login</h2>
      {error && <p style={{ color: 'crimson' }}>{error}</p>}

      <form onSubmit={onSubmit}>
        <div>
          <label htmlFor="login-email">Email</label><br />
          <input
            id="login-email"
            name="email"
            type="email"
            autoComplete="username"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            data-testid="login-email"
          />
        </div>
        <div>
          <label htmlFor="login-password">Password</label><br />
          <input
            id="login-password"
            name="password"
            type="password"
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            data-testid="login-password"
          />
        </div>
        <button type="submit" data-testid="login-submit">
          Login
        </button>
      </form>

      <p>
        No account? <Link to="/register">Register</Link>
      </p>
    </div>
  )
}
