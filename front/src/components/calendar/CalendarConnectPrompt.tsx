import { useState } from 'react'

interface CalendarConnectPromptProps {
  onConnect: () => Promise<void> | void
}

export default function CalendarConnectPrompt({ onConnect }: CalendarConnectPromptProps) {
  const [isLoading, setIsLoading] = useState(false)

  async function handleConnect() {
    setIsLoading(true)
    try {
      await onConnect()
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <section className="calendar-connect-shell">
      <div className="calendar-connect-card">
        <div className="calendar-connect-icon" aria-hidden="true">
          <svg viewBox="0 0 64 64" width="64" height="64">
            <rect x="8" y="10" width="48" height="46" rx="12" fill="#f5f0e8" opacity="0.1" />
            <rect x="12" y="14" width="40" height="38" rx="9" fill="#242220" stroke="rgba(160, 154, 144, 0.18)" />
            <rect x="12" y="14" width="40" height="10" rx="9" fill="#4a6fa5" />
            <circle cx="22" cy="32" r="3.2" fill="#e8a445" />
            <circle cx="32" cy="32" r="3.2" fill="#8fa58b" />
            <circle cx="42" cy="32" r="3.2" fill="#c4826d" />
            <circle cx="22" cy="42" r="3.2" fill="#7eb8b3" />
            <circle cx="32" cy="42" r="3.2" fill="#d4a574" />
          </svg>
        </div>

        <h1 className="calendar-connect-title">Connect your Google Calendar</h1>
        <p className="calendar-connect-subtitle">
          Unlock real time-spending insights from your calendar events.
        </p>

        <button
          type="button"
          className="calendar-connect-cta"
          onClick={handleConnect}
          disabled={isLoading}
          aria-label={isLoading ? 'Connecting Google Calendar' : 'Connect Google Calendar'}
        >
          <span className="calendar-connect-gmark" aria-hidden="true">G</span>
          <span>{isLoading ? 'Connecting...' : 'Connect Google Calendar'}</span>
        </button>

        <p className="calendar-connect-note">Read-only access · Secured via OAuth 2.0</p>
      </div>
    </section>
  )
}
