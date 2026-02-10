import type { ReactElement } from 'react'
import { Navigate } from 'react-router-dom'
import { isAuthenticated } from '@/lib/session'

interface AnonymousOnlyRouteProps {
  children: ReactElement
}

export default function AnonymousOnlyRoute({ children }: AnonymousOnlyRouteProps) {
  if (isAuthenticated()) {
    return <Navigate to="/timespending" replace />
  }

  return children
}
