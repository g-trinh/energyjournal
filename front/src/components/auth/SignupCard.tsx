import { useState, type FormEvent } from 'react'
import { createUser } from '@/services/auth'
import { cn } from '@/lib/utils'

export default function SignupCard() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [feedback, setFeedback] = useState<{
    type: 'success' | 'error'
    message: string
  } | null>(null)

  const [touched, setTouched] = useState<Record<string, boolean>>({})

  const emailValid = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
  const passwordsMatch = password === confirmPassword
  const formValid =
    email.length > 0 &&
    password.length > 0 &&
    confirmPassword.length > 0 &&
    emailValid &&
    passwordsMatch

  function getFieldError(field: string): string | null {
    if (!touched[field]) return null
    if (field === 'email') {
      if (!email) return 'Email is required.'
      if (!emailValid) return 'Enter a valid email address.'
    }
    if (field === 'password' && !password) return 'Password is required.'
    if (field === 'confirmPassword') {
      if (!confirmPassword) return 'Please confirm your password.'
      if (!passwordsMatch) return 'Passwords do not match.'
    }
    return null
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    if (!formValid || submitting) return

    setFeedback(null)
    setSubmitting(true)

    const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone

    const result = await createUser({
      email,
      password,
      confirmPassword,
      timezone,
    })

    setSubmitting(false)

    if (result.ok) {
      setFeedback({
        type: 'success',
        message: 'Check your email to activate your account.',
      })
      setEmail('')
      setPassword('')
      setConfirmPassword('')
      setTouched({})
    } else {
      setFeedback({
        type: 'error',
        message: 'Unable to create account. Please try again.',
      })
    }
  }

  const emailError = getFieldError('email')
  const passwordError = getFieldError('password')
  const confirmError = getFieldError('confirmPassword')

  return (
    <section className="auth-card" aria-labelledby="signup-title">
      <h2 id="signup-title" className="auth-card-title">Sign Up</h2>
      <p className="auth-card-description">Create your account in seconds</p>

      <form className="auth-form" onSubmit={handleSubmit} noValidate>
        <div className="auth-form-fields">
          <div className={cn('auth-field', emailError && 'auth-field-error')}>
            <label htmlFor="signup-email">Email</label>
            <input
              id="signup-email"
              type="email"
              autoComplete="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              onBlur={() => setTouched((t) => ({ ...t, email: true }))}
              aria-invalid={emailError ? 'true' : undefined}
              aria-describedby={emailError ? 'signup-email-error' : undefined}
            />
            {emailError && (
              <p id="signup-email-error" className="auth-field-error-text" role="alert">
                {emailError}
              </p>
            )}
          </div>

          <div className={cn('auth-field', passwordError && 'auth-field-error')}>
            <label htmlFor="signup-password">Password</label>
            <input
              id="signup-password"
              type="password"
              autoComplete="new-password"
              placeholder="••••••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              onBlur={() => setTouched((t) => ({ ...t, password: true }))}
              aria-invalid={passwordError ? 'true' : undefined}
              aria-describedby={passwordError ? 'signup-password-error' : undefined}
            />
            {passwordError && (
              <p id="signup-password-error" className="auth-field-error-text" role="alert">
                {passwordError}
              </p>
            )}
          </div>

          <div className={cn('auth-field', confirmError && 'auth-field-error')}>
            <label htmlFor="signup-confirm-password">Confirm Password</label>
            <input
              id="signup-confirm-password"
              type="password"
              autoComplete="new-password"
              placeholder="••••••••••••"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              onBlur={() =>
                setTouched((t) => ({ ...t, confirmPassword: true }))
              }
              aria-invalid={confirmError ? 'true' : undefined}
              aria-describedby={
                confirmError ? 'signup-confirm-error' : undefined
              }
            />
            {confirmError && (
              <p id="signup-confirm-error" className="auth-field-error-text" role="alert">
                {confirmError}
              </p>
            )}
          </div>
        </div>

        <div className="auth-cta-area">
          {feedback && (
            <div
              className={cn(
                'auth-feedback',
                feedback.type === 'success'
                  ? 'auth-feedback-success'
                  : 'auth-feedback-error',
              )}
              role="status"
              aria-live="polite"
            >
              {feedback.message}
            </div>
          )}

          <button
            type="submit"
            className="auth-btn auth-btn-signup"
            disabled={!formValid || submitting}
          >
            {submitting ? 'Creating account…' : 'Create Account'}
          </button>
        </div>
      </form>
    </section>
  )
}
