import { useEffect, useState, useRef } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { activateUser } from '@/services/auth'
import '../styles/auth.css'

type Status = 'idle' | 'loading' | 'success' | 'error'

export default function ActivatePage() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const token = searchParams.get('token')

  const [status, setStatus] = useState<Status>(token ? 'loading' : 'error')
  const calledRef = useRef(false)

  useEffect(() => {
    if (!token || calledRef.current) return
    calledRef.current = true

    activateUser(token).then((result) => {
      if (result.ok) {
        setStatus('success')
      } else {
        setStatus('error')
      }
    })
  }, [token])

  useEffect(() => {
    if (status !== 'success') return
    const timer = setTimeout(() => navigate('/auth', { replace: true }), 5000)
    return () => clearTimeout(timer)
  }, [status, navigate])

  return (
    <div className="app">
      <div className="ambient-glow ambient-glow-1" />
      <div className="ambient-glow ambient-glow-2" />
      <div className="grain-overlay" />

      <main className="activate-content">
        <div className="activate-card">
          <h1 className="auth-headline">Account Activation</h1>

          <div className="activate-status" role="status" aria-live="polite">
            {status === 'loading' && (
              <>
                <div className="activate-spinner" aria-hidden="true">
                  <div className="spinner-ring" />
                  <div className="spinner-ring" />
                </div>
                <p className="auth-card-description">
                  Activating your account…
                </p>
              </>
            )}

            {status === 'success' && (
              <>
                <div className="auth-feedback auth-feedback-success" style={{ display: 'inline-block' }}>
                  Your account has been activated successfully.
                </div>
                <p className="activate-redirect-note">
                  Redirecting to login…
                </p>
              </>
            )}

            {status === 'error' && (
              <div className="auth-feedback auth-feedback-error" style={{ display: 'inline-block' }}>
                {!token
                  ? 'This activation link is invalid or has expired.'
                  : 'Unable to activate your account. The link may be invalid or has expired.'}
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  )
}
