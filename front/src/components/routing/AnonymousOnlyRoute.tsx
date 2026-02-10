import { trackEvent } from '@/lib/analytics'
import { isAuthenticated } from '@/lib/session'
import type { ReactElement } from 'react'
import { Navigate, useLocation } from 'react-router-dom'

interface AnonymousOnlyRouteProps {
  children: ReactElement
}

export default function AnonymousOnlyRoute({ children }: AnonymousOnlyRouteProps) {
  const location = useLocation()

  if (isAuthenticated()) {
    trackEvent('auth_redirect_blocked', {
      guard: 'anonymous-only',
      from: location.pathname,
      to: '/timespending',
    })

    return <Navigate to="/timespending" replace />
  }

  return children
}
