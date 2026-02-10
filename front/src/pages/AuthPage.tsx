import LoginCard from '@/components/auth/LoginCard'
import SignupCard from '@/components/auth/SignupCard'
import { persistSession } from '@/lib/session'
import type { AuthTokensResponse } from '@/services/auth'
import { useNavigate } from 'react-router-dom'
import '../styles/auth.css'

export default function AuthPage() {
  const navigate = useNavigate()

  function handleLoginSuccess(tokens: AuthTokensResponse) {
    persistSession(tokens.idToken, tokens.refreshToken)
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
