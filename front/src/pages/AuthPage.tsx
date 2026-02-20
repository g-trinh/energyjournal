import LoginCard from '@/components/auth/LoginCard'
import { trackEvent } from '@/lib/analytics'
import SignupCard from '@/components/auth/SignupCard'
import { useAuth } from '@/contexts/AuthContext'
import { persistSession } from '@/lib/session'
import type { AuthTokensResponse } from '@/services/auth'
import { useNavigate } from 'react-router-dom'
import '../styles/auth.css'

export default function AuthPage() {
  const navigate = useNavigate()
  const { signIn } = useAuth()

  function handleLoginSuccess(tokens: AuthTokensResponse, email: string) {
    persistSession(tokens.idToken, tokens.refreshToken)
    signIn(email)
    trackEvent('auth_login_success')
    navigate('/timespending', { replace: true })
  }

  return (
    <div className="app">
      <div className="ambient-glow ambient-glow-1" />
      <div className="ambient-glow ambient-glow-2" />
      <div className="grain-overlay" />

      <main className="auth-content">
        <header className="auth-header">
          <h1 className="auth-headline">Access your account</h1>
          <p className="auth-subline">
            Access your Energy Journal or create a new account in one place.
          </p>
        </header>

        <div className="auth-cards">
          <LoginCard onLoginSuccess={handleLoginSuccess} />

          <SignupCard />
        </div>
      </main>
    </div>
  )
}
