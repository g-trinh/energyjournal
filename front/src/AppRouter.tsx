import { Navigate, Route, Routes } from 'react-router-dom'
import App from './App'
import Layout from './components/layout/Layout'
import AnonymousOnlyRoute from './components/routing/AnonymousOnlyRoute'
import ProtectedRoute from './components/routing/ProtectedRoute'
import ActivatePage from './pages/ActivatePage'
import AuthPage from './pages/AuthPage'
import EnergyLevelsEditPage from './pages/EnergyLevelsEditPage'
import LandingPage from './pages/LandingPage'

export default function AppRouter() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<LandingPage />} />
        <Route
          path="/timespending"
          element={
            <ProtectedRoute>
              <App />
            </ProtectedRoute>
          }
        />
        <Route
          path="/energy/levels/edit"
          element={
            <ProtectedRoute>
              <EnergyLevelsEditPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/auth"
          element={
            <AnonymousOnlyRoute>
              <AuthPage />
            </AnonymousOnlyRoute>
          }
        />
        <Route
          path="/activate"
          element={
            <AnonymousOnlyRoute>
              <ActivatePage />
            </AnonymousOnlyRoute>
          }
        />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  )
}
