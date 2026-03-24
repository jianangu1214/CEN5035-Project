import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { getFlights, deleteFlight } from '../api/flights'

function formatDateTime(str) {
  if (!str) return '—'
  const d = new Date(str)
  return isNaN(d.getTime()) ? str : d.toLocaleString()
}

function formatPrice(n) {
  if (n == null || n === '') return '—'
  const num = Number(n)
  return isNaN(num) ? n : `$${num.toFixed(2)}`
}

export default function FlightList() {
  const [flights, setFlights] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [deletingId, setDeletingId] = useState(null)

  const fetchFlights = async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await getFlights()
      setFlights(Array.isArray(data) ? data : [])
    } catch (err) {
      setError(err.body?.error || err.body?.message || err.message || 'Failed to load flights')
      setFlights([])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchFlights()
  }, [])

  const handleDelete = async (id, label) => {
    if (!window.confirm(`Delete flight "${label || 'this entry'}"?`)) return
    setDeletingId(id)
    try {
      await deleteFlight(id)
      setFlights((prev) => prev.filter((f) => f.id !== id))
    } catch (err) {
      setError(err.body?.error || err.body?.message || err.message || 'Failed to delete flight')
    } finally {
      setDeletingId(null)
    }
  }

  if (loading) {
    return <div className="loading">Loading flights…</div>
  }

  return (
    <>
      <div className="page-header">
        <h1>Flights</h1>
        <Link to="/flights/new" className="btn btn-primary">
          Add flight
        </Link>
      </div>

      {error && (
        <div className="error-banner" role="alert">
          {error}
        </div>
      )}

      <div className="table-wrap">
        {flights.length === 0 ? (
          <div className="empty-state">
            <p>No flights yet.</p>
            <Link to="/flights/new" className="btn btn-primary">
              Add your first flight
            </Link>
          </div>
        ) : (
          <table className="table">
            <thead>
              <tr>
                <th>Airline</th>
                <th>Flight</th>
                <th>From</th>
                <th>To</th>
                <th>Depart</th>
                <th>Arrive</th>
                <th>Price</th>
                <th>Notes</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {flights.map((f) => (
                <tr key={f.id}>
                  <td>{f.airline || '—'}</td>
                  <td>{f.flight_number || '—'}</td>
                  <td>{f.from_airport || '—'}</td>
                  <td>{f.to_airport || '—'}</td>
                  <td>{formatDateTime(f.depart_time)}</td>
                  <td>{formatDateTime(f.arrive_time)}</td>
                  <td>{formatPrice(f.price)}</td>
                  <td>{f.notes ? String(f.notes).slice(0, 40) + (String(f.notes).length > 40 ? '…' : '') : '—'}</td>
                  <td>
                    <span className="actions">
                      <Link to={`/flights/${f.id}/edit`} className="btn btn-secondary">
                        Edit
                      </Link>
                      <button
                        type="button"
                        className="btn btn-danger"
                        disabled={deletingId === f.id}
                        onClick={() => handleDelete(f.id, `${f.airline} ${f.flight_number}`)}
                      >
                        {deletingId === f.id ? 'Deleting…' : 'Delete'}
                      </button>
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </>
  )
}
