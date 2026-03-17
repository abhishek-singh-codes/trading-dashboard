import { Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import LoginPage     from './pages/LoginPage'
import RegisterPage  from './pages/RegisterPage'
import DashboardPage from './pages/DashboardPage'
import StockPage     from './pages/StockPage'

function PrivateRoute({ children }) {
  const { isAuthenticated } = useAuth()
  return isAuthenticated ? children : <Navigate to="/login" replace />
}

export default function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route path="/login"         element={<LoginPage />} />
        <Route path="/register"      element={<RegisterPage />} />
        <Route path="/"              element={<PrivateRoute><DashboardPage /></PrivateRoute>} />
        <Route path="/stock/:symbol" element={<PrivateRoute><StockPage /></PrivateRoute>} />
        <Route path="*"              element={<Navigate to="/" replace />} />
      </Routes>
    </AuthProvider>
  )
}