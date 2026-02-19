import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { request } from '../api/client'

export default function Register() {
  const nav = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [ok, setOk] = useState(false)

  async function onSubmit(e) {
    e.preventDefault()
    setError('')
    setOk(false)
    try {
      await request('/auth/register', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      })
      setOk(true)
      setTimeout(() => nav('/login'), 300)
    } catch (err) {
      setError(err?.message || 'Register failed')
    }
  }

  return (
    <div>
      <h2>Register</h2>
      {error && <p style={{ color: 'crimson' }}>{error}</p>}
      {ok && <p style={{ color: 'green' }}>Registered! Redirecting to login…</p>}

      <form onSubmit={onSubmit}>
        <div>
          <label>Email</label><br />
          <input value={email} onChange={(e) => setEmail(e.target.value)} />
        </div>
        <div>
          <label>Password</label><br />
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        </div>
        <button type="submit">Create account</button>
      </form>

      <p>
        Already have an account? <Link to="/login">Login</Link>
      </p>
    </div>
  )
}
