export type AnalyticsEventName =
  | 'auth_login_success'
  | 'auth_redirect_blocked'
  | 'timespending_load_failed'

export type AnalyticsPayload = Record<string, string | number | boolean | null>

export function trackEvent(name: AnalyticsEventName, payload: AnalyticsPayload = {}): void {
  const event = {
    name,
    payload,
    timestamp: new Date().toISOString(),
  }

  window.dispatchEvent(new CustomEvent('analytics:event', { detail: event }))
}
