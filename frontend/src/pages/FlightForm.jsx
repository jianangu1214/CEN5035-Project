import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { getFlight, createFlight, updateFlight } from '../api/flights'

const emptyForm = {
  airline: '',
  flight_number: '',
  from_airport: '',
  to_airport: '',
  depart_time: '',
  arrive_time: '',
  price: '',
  notes: '',
}

/** datetime-local value from RFC3339 */
function toInputDateTime(value) {
  if (!value) return ''
  const d = new Date(value)
  if (isNaN(d.getTime())) return ''
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function toRFC3339(localValue) {
  if (!localValue) return null
  const d = new Date(localValue)
  return isNaN(d.getTime()) ? null : d.toISOString()
}

export default function FlightForm() {
  const { id } = useParams()
  const navigate = useNavigate()
  const isEdit = Boolean(id)

  const [form, setForm] = useState(emptyForm)
  const [loading, setLoading] = useState(isEdit)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState(null)

  useEffect(() => {
    if (!isEdit) return
    let cancelled = false
    setLoading(true)
    setError(null)
    getFlight(id)
      .then((data) => {
        if (!cancelled) {
          setForm({
            airline: data.airline ?? '',
            flight_number: data.flight_number ?? '',
            from_airport: (data.from_airport ?? '').slice(0, 3).toUpperCase(),
            to_airport: (data.to_airport ?? '').slice(0, 3).toUpperCase(),
            depart_time: toInputDateTime(data.depart_time),
            arrive_time: toInputDateTime(data.arrive_time),
            price: data.price != null ? String(data.price) : '',
            notes: data.notes ?? '',
          })
        }
      })
      .catch((err) => {
        if (!cancelled) setError(err.message || 'Failed to load flight')
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => { cancelled = true }
  }, [id, isEdit])

  const handleChange = (e) => {
    const { name, value } = e.target
    if (name === 'from_airport' || name === 'to_airport') {
      setForm((prev) => ({ ...prev, [name]: value.toUpperCase().slice(0, 3) }))
      return
    }
    setForm((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError(null)
    setSaving(true)
    const depart = toRFC3339(form.depart_time)
    const arrive = toRFC3339(form.arrive_time)
    const payload = {
      airline: form.airline.trim(),
      flight_number: form.flight_number.trim(),
      from_airport: form.from_airport.trim().toUpperCase(),
      to_airport: form.to_airport.trim().toUpperCase(),
      depart_time: depart,
      arrive_time: arrive,
      price: form.price === '' ? 0 : Number(form.price),
      notes: form.notes.trim() || '',
    }
    try {
      if (!payload.airline || !payload.flight_number || !payload.from_airport || !payload.to_airport) {
        throw new Error('Airline, flight number, and airports are required')
      }
      if (payload.from_airport.length !== 3 || payload.to_airport.length !== 3) {
        throw new Error('Airport codes must be 3 letters')
      }
      if (!depart || !arrive) {
        throw new Error('Depart and arrive times are required')
      }
      if (isEdit) {
        await updateFlight(id, payload)
        navigate('/flights')
      } else {
        await createFlight(payload)
        navigate('/flights')
      }
    } catch (err) {
      setError(err.message || (isEdit ? 'Failed to update flight' : 'Failed to add flight'))
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return <div className="loading">Loading…</div>
  }

  return (
    <>
      <div className="page-header">
        <h1>{isEdit ? 'Edit flight' : 'Add flight'}</h1>
      </div>

      {error && (
        <div className="error-banner" role="alert">
          {error}
        </div>
      )}

      <div className="form-card">
        <form onSubmit={handleSubmit} data-testid="flight-form">
          <div className="form-group">
            <label htmlFor="airline">Airline</label>
            <input
              id="airline"
              name="airline"
              value={form.airline}
              onChange={handleChange}
              placeholder="e.g. Delta"
              data-testid="flight-airline"
            />
          </div>
          <div className="form-group">
            <label htmlFor="flight_number">Flight number</label>
            <input
              id="flight_number"
              name="flight_number"
              value={form.flight_number}
              onChange={handleChange}
              placeholder="e.g. DL123"
              data-testid="flight-number"
            />
          </div>
          <div className="form-group">
            <label htmlFor="from_airport">From (IATA)</label>
            <input
              id="from_airport"
              name="from_airport"
              value={form.from_airport}
              onChange={handleChange}
              maxLength={3}
              placeholder="MIA"
              data-testid="flight-from"
            />
          </div>
          <div className="form-group">
            <label htmlFor="to_airport">To (IATA)</label>
            <input
              id="to_airport"
              name="to_airport"
              value={form.to_airport}
              onChange={handleChange}
              maxLength={3}
              placeholder="JFK"
              data-testid="flight-to"
            />
          </div>
          <div className="form-group">
            <label htmlFor="depart_time">Depart</label>
            <input
              id="depart_time"
              name="depart_time"
              type="datetime-local"
              value={form.depart_time}
              onChange={handleChange}
              data-testid="flight-depart"
            />
          </div>
          <div className="form-group">
            <label htmlFor="arrive_time">Arrive</label>
            <input
              id="arrive_time"
              name="arrive_time"
              type="datetime-local"
              value={form.arrive_time}
              onChange={handleChange}
              data-testid="flight-arrive"
            />
          </div>
          <div className="form-group">
            <label htmlFor="price">Price ($)</label>
            <input
              id="price"
              name="price"
              type="number"
              min="0"
              step="0.01"
              value={form.price}
              onChange={handleChange}
              placeholder="0.00"
              data-testid="flight-price"
            />
          </div>
          <div className="form-group">
            <label htmlFor="notes">Notes (optional)</label>
            <textarea
              id="notes"
              name="notes"
              value={form.notes}
              onChange={handleChange}
              placeholder="Any extra details…"
            />
          </div>
          <div className="form-actions">
            <button type="submit" className="btn btn-primary" disabled={saving} data-testid="flight-submit">
              {saving ? 'Saving…' : isEdit ? 'Update' : 'Add flight'}
            </button>
            <button
              type="button"
              className="btn btn-secondary"
              onClick={() => navigate('/flights')}
              disabled={saving}
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </>
  )
}
