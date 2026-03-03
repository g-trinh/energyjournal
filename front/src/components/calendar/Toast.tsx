import { useEffect } from 'react'

interface ToastProps {
  message: string
  subtitle?: string
  duration?: number
  onDismiss: () => void
}

export default function Toast({ message, subtitle, duration = 4000, onDismiss }: ToastProps) {
  useEffect(() => {
    const timeout = window.setTimeout(onDismiss, duration)
    return () => window.clearTimeout(timeout)
  }, [duration, onDismiss])

  return (
    <aside className="calendar-toast" role="status" aria-live="polite">
      <span className="calendar-toast-accent" aria-hidden="true" />
      <span className="calendar-toast-check" aria-hidden="true">✓</span>
      <div className="calendar-toast-copy">
        <p className="calendar-toast-title">{message}</p>
        {subtitle && <p className="calendar-toast-subtitle">{subtitle}</p>}
      </div>
    </aside>
  )
}
