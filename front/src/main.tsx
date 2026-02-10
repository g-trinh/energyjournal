import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import App from './App.tsx'
import './index.css'
import AnonymousOnlyRoute from './components/routing/AnonymousOnlyRoute.tsx'
import ProtectedRoute from './components/routing/ProtectedRoute.tsx'
import ActivatePage from './pages/ActivatePage.tsx'
import AuthPage from './pages/AuthPage.tsx'
import LandingPage from './pages/LandingPage.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
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
      </Routes>
    </BrowserRouter>
  </StrictMode>,
)
