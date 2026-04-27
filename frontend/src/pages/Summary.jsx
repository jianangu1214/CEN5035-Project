import { useEffect, useState } from 'react'
import { getSummary } from '../api/summary'

export default function Summary() {
  const [type, setType] = useState('month')
  const [data, setData] = useState([])

  useEffect(() => {
    getSummary(type).then(setData)
  }, [type])

  return (
    <div>
      <h1>Summary</h1>

      <select value={type} onChange={(e) => setType(e.target.value)}>
        <option value="month">Month</option>
        <option value="quarter">Quarter</option>
        <option value="year">Year</option>
      </select>

      <table>
        <thead>
          <tr>
            <th>Period</th>
            <th>Flights</th>
            <th>Hotels</th>
            <th>Nights</th>
            <th>Spend</th>
          </tr>
        </thead>
        <tbody>
          {data.map((row, i) => (
            <tr key={i}>
              <td>{row.period}</td>
              <td>{row.flights}</td>
              <td>{row.hotels}</td>
              <td>{row.nights}</td>
              <td>{row.spend}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}