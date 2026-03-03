import { getIdToken } from '@/lib/session'
import type { CalendarItem } from '@/components/calendar/types'

export type CalendarStatus = 'disconnected' | 'pending_selection' | 'connected'

function authHeaders() {
  const token = getIdToken()
  if (!token) {
    return null
  }
  return { Authorization: `Bearer ${token}` }
}

export async function getCalendarStatus(): Promise<CalendarStatus> {
  const headers = authHeaders()
  if (!headers) {
    throw new Error('unauthorized')
  }

  const response = await fetch('/api/calendar/status', { headers })
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`)
  }
  const payload = await response.json() as { status: CalendarStatus }
  return payload.status
}

export async function getCalendarAuthURL(): Promise<string> {
  const headers = authHeaders()
  if (!headers) {
    throw new Error('unauthorized')
  }

  const response = await fetch('/api/calendar/auth', { headers })
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`)
  }
  const payload = await response.json() as { auth_url: string }
  return payload.auth_url
}

export async function getCalendars(): Promise<CalendarItem[]> {
  const headers = authHeaders()
  if (!headers) {
    throw new Error('unauthorized')
  }

  const response = await fetch('/api/calendar/calendars', { headers })
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`)
  }
  return response.json() as Promise<CalendarItem[]>
}

export async function saveCalendarConnection(calendarID: string): Promise<void> {
  const headers = authHeaders()
  if (!headers) {
    throw new Error('unauthorized')
  }

  const response = await fetch('/api/calendar/connection', {
    method: 'PUT',
    headers: {
      ...headers,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ calendar_id: calendarID }),
  })
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`)
  }
}
