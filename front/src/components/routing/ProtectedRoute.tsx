import { trackEvent } from '@/lib/analytics'
import { isAuthenticated } from '@/lib/session'
import type { ReactElement } from 'react'
import { Navigate, useLocation } from 'react-router-dom'

interface ProtectedRouteProps {
  children: ReactElement
}

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  const location = useLocation()

  if (!isAuthenticated()) {
    trackEvent('auth_redirect_blocked', {
      guard: 'protected',
      from: location.pathname,
      to: '/auth',
    })

    return <Navigate to="/auth" replace />
  }

  return children
}
