import { Routes, Route, Navigate } from 'react-router-dom'
import NavBar from './components/NavBar'
import PrivateRoute from './components/PrivateRoute'
import HotelList from './pages/HotelList'
import HotelForm from './pages/HotelForm'
import FlightList from './pages/FlightList'
import FlightForm from './pages/FlightForm'
import Login from './pages/Login'
import Register from './pages/Register'
import './App.css'

export default function App() {
  return (
    <div className="app">
      <NavBar />
      <main className="main">
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route element={<PrivateRoute />}>
            <Route path="/" element={<Navigate to="/hotels" replace />} />
            <Route path="/hotels" element={<HotelList />} />
            <Route path="/hotels/new" element={<HotelForm />} />
            <Route path="/hotels/:id/edit" element={<HotelForm />} />
            <Route path="/flights" element={<FlightList />} />
            <Route path="/flights/new" element={<FlightForm />} />
            <Route path="/flights/:id/edit" element={<FlightForm />} />
          </Route>
          <Route path="*" element={<Navigate to="/hotels" replace />} />
        </Routes>
      </main>
    </div>
  )
}
