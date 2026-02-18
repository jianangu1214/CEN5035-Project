import { Routes, Route, Link } from 'react-router-dom'
import HotelList from './pages/HotelList'
import HotelForm from './pages/HotelForm'
import './App.css'

function App() {
  return (
    <div className="app">
      <nav className="nav">
        <Link to="/">TravelLog</Link>
        <Link to="/hotels">Hotels</Link>
        <Link to="/hotels/new">Add Hotel</Link>
      </nav>
      <main className="main">
        <Routes>
          <Route path="/" element={<HotelList />} />
          <Route path="/hotels" element={<HotelList />} />
          <Route path="/hotels/new" element={<HotelForm />} />
          <Route path="/hotels/:id/edit" element={<HotelForm />} />
        </Routes>
      </main>
    </div>
  )
}

export default App
