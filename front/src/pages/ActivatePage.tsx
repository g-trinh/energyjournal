import { useSearchParams } from 'react-router-dom'
import '../styles/auth.css'

export default function ActivatePage() {
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token')

  return (
    <div className="app">
      <div className="ambient-glow ambient-glow-1" />
      <div className="ambient-glow ambient-glow-2" />
      <div className="grain-overlay" />

      <main className="activate-content">
        <div className="activate-card">
          <h1 className="auth-headline">Account Activation</h1>
          {!token ? (
            <p className="auth-card-description">
              This activation link is invalid or has expired.
            </p>
          ) : (
            <p className="auth-card-description">Activating your account…</p>
          )}
          {/* Full activation logic will be implemented in AUTH-005 */}
        </div>
      </main>
    </div>
  )
}
