import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { getHotels, deleteHotel } from '../api/hotels'

function formatDate(str) {
  if (!str) return '—'
  const d = new Date(str)
  return isNaN(d.getTime()) ? str : d.toLocaleDateString()
}

function formatPrice(n) {
  if (n == null || n === '') return '—'
  const num = Number(n)
  return isNaN(num) ? n : `$${num.toFixed(2)}`
}

export default function HotelList() {
  const [hotels, setHotels] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [deletingId, setDeletingId] = useState(null)

  const fetchHotels = async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await getHotels()
      setHotels(Array.isArray(data) ? data : [])
    } catch (err) {
      setError(err.body?.error || err.body?.message || err.message || 'Failed to load hotels')
      setHotels([])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchHotels()
  }, [])

  const handleDelete = async (id, hotelName) => {
    if (!window.confirm(`Delete "${hotelName || 'this hotel'}"?`)) return
    setDeletingId(id)
    try {
      await deleteHotel(id)
      setHotels((prev) => prev.filter((h) => h.id !== id))
    } catch (err) {
      setError(err.body?.error || err.body?.message || err.message || 'Failed to delete hotel')
    } finally {
      setDeletingId(null)
    }
  }

  if (loading) {
    return <div className="loading">Loading hotels…</div>
  }

  return (
    <>
      <div className="page-header">
        <h1>Hotel stays</h1>
        <Link to="/hotels/new" className="btn btn-primary">
          Add hotel
        </Link>
      </div>

      {error && (
        <div className="error-banner" role="alert">
          {error}
        </div>
      )}

      <div className="table-wrap">
        {hotels.length === 0 ? (
          <div className="empty-state">
            <p>No hotel stays yet.</p>
            <Link to="/hotels/new" className="btn btn-primary">
              Add your first hotel
            </Link>
          </div>
        ) : (
          <table className="table">
            <thead>
              <tr>
                <th>Hotel</th>
                <th>City</th>
                <th>Country</th>
                <th>Check-in</th>
                <th>Check-out</th>
                <th>Price</th>
                <th>Notes</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {hotels.map((h) => (
                <tr key={h.id}>
                  <td>{h.hotel_name || '—'}</td>
                  <td>{h.city || '—'}</td>
                  <td>{h.country || '—'}</td>
                  <td>{formatDate(h.check_in)}</td>
                  <td>{formatDate(h.check_out)}</td>
                  <td>{formatPrice(h.price)}</td>
                  <td>{h.notes ? String(h.notes).slice(0, 40) + (String(h.notes).length > 40 ? '…' : '') : '—'}</td>
                  <td>
                    <span className="actions">
                      <Link to={`/hotels/${h.id}/edit`} className="btn btn-secondary">
                        Edit
                      </Link>
                      <button
                        type="button"
                        className="btn btn-danger"
                        disabled={deletingId === h.id}
                        onClick={() => handleDelete(h.id, h.hotel_name)}
                      >
                        {deletingId === h.id ? 'Deleting…' : 'Delete'}
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
