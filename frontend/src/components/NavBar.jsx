import { Link, useLocation, useNavigate } from 'react-router-dom'
import { getAuthToken, setAuthToken } from '../api/client'

export default function NavBar() {
  const location = useLocation()
  const navigate = useNavigate()
  const token = getAuthToken()

  function logout() {
    setAuthToken(null)
    navigate('/login')
  }

  return (
    <nav className="nav" data-path={location.pathname}>
      <Link to="/">TravelLog</Link>
      {token && (
        <>
          <Link to="/hotels">Hotels</Link>
          <Link to="/hotels/new">Add hotel</Link>
          <Link to="/flights">Flights</Link>
          <Link to="/flights/new">Add flight</Link>
        </>
      )}
      {token ? (
        <button type="button" className="btn btn-secondary" onClick={logout}>
          Logout
        </button>
      ) : (
        <>
          <Link to="/login">Login</Link>
          <Link to="/register">Register</Link>
        </>
      )}
    </nav>
  )
}
