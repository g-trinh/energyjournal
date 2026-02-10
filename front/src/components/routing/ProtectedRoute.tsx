import type { ReactElement } from 'react'
import { Navigate } from 'react-router-dom'
import { isAuthenticated } from '@/lib/session'

interface ProtectedRouteProps {
  children: ReactElement
}

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  if (!isAuthenticated()) {
    return <Navigate to="/auth" replace />
  }

  return children
}
