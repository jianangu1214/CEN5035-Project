import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { getHotel, createHotel, updateHotel } from '../api/hotels'

const emptyForm = {
  hotel_name: '',
  city: '',
  country: '',
  check_in: '',
  check_out: '',
  price: '',
  notes: '',
}

function toInputDate(value) {
  if (!value) return ''
  const d = new Date(value)
  return isNaN(d.getTime()) ? '' : d.toISOString().slice(0, 10)
}

export default function HotelForm() {
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
    getHotel(id)
      .then((data) => {
        if (!cancelled) {
          setForm({
            hotel_name: data.hotel_name ?? '',
            city: data.city ?? '',
            country: data.country ?? '',
            check_in: toInputDate(data.check_in),
            check_out: toInputDate(data.check_out),
            price: data.price != null ? String(data.price) : '',
            notes: data.notes ?? '',
          })
        }
      })
      .catch((err) => {
        if (!cancelled) setError(err.message || 'Failed to load hotel')
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => { cancelled = true }
  }, [id, isEdit])

  const handleChange = (e) => {
    const { name, value } = e.target
    setForm((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError(null)
    setSaving(true)
    const payload = {
      hotel_name: form.hotel_name.trim() || null,
      city: form.city.trim() || null,
      country: form.country.trim() || null,
      check_in: form.check_in || null,
      check_out: form.check_out || null,
      price: form.price === '' ? null : Number(form.price),
      notes: form.notes.trim() || null,
    }
    try {
      if (isEdit) {
        await updateHotel(id, payload)
        navigate('/hotels')
      } else {
        await createHotel(payload)
        navigate('/hotels')
      }
    } catch (err) {
      setError(err.message || (isEdit ? 'Failed to update hotel' : 'Failed to add hotel'))
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
        <h1>{isEdit ? 'Edit hotel stay' : 'Add hotel stay'}</h1>
      </div>

      {error && (
        <div className="error-banner" role="alert">
          {error}
        </div>
      )}

      <div className="form-card">
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="hotel_name">Hotel name</label>
            <input
              id="hotel_name"
              name="hotel_name"
              value={form.hotel_name}
              onChange={handleChange}
              placeholder="e.g. Grand Hotel"
            />
          </div>
          <div className="form-group">
            <label htmlFor="city">City</label>
            <input
              id="city"
              name="city"
              value={form.city}
              onChange={handleChange}
              placeholder="e.g. Miami"
            />
          </div>
          <div className="form-group">
            <label htmlFor="country">Country</label>
            <input
              id="country"
              name="country"
              value={form.country}
              onChange={handleChange}
              placeholder="e.g. USA"
            />
          </div>
          <div className="form-group">
            <label htmlFor="check_in">Check-in date</label>
            <input
              id="check_in"
              name="check_in"
              type="date"
              value={form.check_in}
              onChange={handleChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="check_out">Check-out date</label>
            <input
              id="check_out"
              name="check_out"
              type="date"
              value={form.check_out}
              onChange={handleChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="price">Total price ($)</label>
            <input
              id="price"
              name="price"
              type="number"
              min="0"
              step="0.01"
              value={form.price}
              onChange={handleChange}
              placeholder="0.00"
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
            <button type="submit" className="btn btn-primary" disabled={saving}>
              {saving ? 'Saving…' : isEdit ? 'Update' : 'Add hotel'}
            </button>
            <button
              type="button"
              className="btn btn-secondary"
              onClick={() => navigate('/hotels')}
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
