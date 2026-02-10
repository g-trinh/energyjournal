import { useState, type FormEvent } from 'react'
import { login, type AuthTokensResponse } from '@/services/auth'
import { cn } from '@/lib/utils'

interface LoginCardProps {
  onLoginSuccess: (tokens: AuthTokensResponse) => void
}

export default function LoginCard({ onLoginSuccess }: LoginCardProps) {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [feedback, setFeedback] = useState<{
    type: 'error'
    message: string
  } | null>(null)

  const [touched, setTouched] = useState<Record<string, boolean>>({})

  const emailValid = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
  const formValid = email.length > 0 && password.length > 0 && emailValid

  function getFieldError(field: string): string | null {
    if (!touched[field]) return null
    if (field === 'email') {
      if (!email) return 'Email is required.'
      if (!emailValid) return 'Enter a valid email address.'
    }
    if (field === 'password' && !password) return 'Password is required.'
    return null
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    if (!formValid || submitting) return

    setFeedback(null)
    setSubmitting(true)

    const result = await login({ email, password })

    setSubmitting(false)

    if (result.ok) {
      onLoginSuccess(result.data)
    } else {
      setFeedback({ type: 'error', message: 'Invalid email or password.' })
    }
  }

  const emailError = getFieldError('email')
  const passwordError = getFieldError('password')

  return (
    <section className="auth-card" aria-labelledby="login-title">
      <h2 id="login-title" className="auth-card-title">Log In</h2>
      <p className="auth-card-description">Return to your dashboard</p>

      <form className="auth-form" onSubmit={handleSubmit} noValidate>
        <div className="auth-form-fields">
          <div className={cn('auth-field', emailError && 'auth-field-error')}>
            <label htmlFor="login-email">Email</label>
            <input
              id="login-email"
              type="email"
              autoComplete="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              onBlur={() => setTouched((t) => ({ ...t, email: true }))}
              aria-invalid={emailError ? 'true' : undefined}
              aria-describedby={emailError ? 'login-email-error' : undefined}
            />
            {emailError && (
              <p id="login-email-error" className="auth-field-error-text" role="alert">
                {emailError}
              </p>
            )}
          </div>

          <div className={cn('auth-field', passwordError && 'auth-field-error')}>
            <label htmlFor="login-password">Password</label>
            <input
              id="login-password"
              type="password"
              autoComplete="current-password"
              placeholder="••••••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              onBlur={() => setTouched((t) => ({ ...t, password: true }))}
              aria-invalid={passwordError ? 'true' : undefined}
              aria-describedby={passwordError ? 'login-password-error' : undefined}
            />
            {passwordError && (
              <p id="login-password-error" className="auth-field-error-text" role="alert">
                {passwordError}
              </p>
            )}
          </div>
        </div>

        <div className="auth-cta-area">
          {feedback && (
            <div
              className="auth-feedback auth-feedback-error"
              role="alert"
              aria-live="polite"
            >
              {feedback.message}
            </div>
          )}

          <button
            type="submit"
            className="auth-btn auth-btn-login"
            disabled={!formValid || submitting}
          >
            {submitting ? 'Signing in…' : 'Log In'}
          </button>
        </div>
      </form>
    </section>
  )
}
