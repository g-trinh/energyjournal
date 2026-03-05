import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from 'recharts'
import { useNavigate } from 'react-router-dom'
import CalendarConnectPrompt from '@/components/calendar/CalendarConnectPrompt'
import CalendarPicker from '@/components/calendar/CalendarPicker'
import Toast from '@/components/calendar/Toast'
import type { CalendarItem } from '@/components/calendar/types'
import { trackEvent } from '@/lib/analytics'
import { clearSession, getIdToken } from '@/lib/session'
import { getSpendingColor, toChartData, truncateAxisLabel, type ChartData, type Spendings } from '@/lib/timeSpending'
import {
  getCalendarAuthURL,
  getCalendars,
  getCalendarStatus,
  saveCalendarConnection,
  type CalendarStatus,
} from '@/services/calendar'
import './App.css'
import './styles/calendar.css'

type FetchErrorKind = 'offline' | 'generic'
type ViewStatus = 'loading' | CalendarStatus

function startOfISOMonday(date: Date): Date {
  const normalized = new Date(date)
  normalized.setHours(0, 0, 0, 0)
  const day = normalized.getDay()
  const offset = day === 0 ? 6 : day - 1
  normalized.setDate(normalized.getDate() - offset)
  return normalized
}

function addDays(date: Date, days: number): Date {
  const next = new Date(date)
  next.setDate(next.getDate() + days)
  return next
}

function formatWeekLabel(weekStart: Date): string {
  const weekEnd = addDays(weekStart, 6)
  const startLabel = weekStart.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  const endLabel = weekEnd.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  const yearLabel = weekEnd.getFullYear()
  return `${startLabel} – ${endLabel}, ${yearLabel}`
}

function formatApiDate(date: Date): string {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// Custom tooltip component
const CustomTooltip = ({ active, payload }: { active?: boolean; payload?: Array<{ value: number; payload: ChartData }> }) => {
  if (active && payload && payload.length) {
    const data = payload[0]
    return (
      <div className="custom-tooltip">
        <p className="tooltip-label">{data.payload.name}</p>
        <p className="tooltip-value">{data.value.toFixed(1)}h</p>
      </div>
    )
  }
  return null
}

function App() {
  const navigate = useNavigate()
  const [data, setData] = useState<ChartData[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [errorKind, setErrorKind] = useState<FetchErrorKind>('generic')
  const [calendarStatus, setCalendarStatus] = useState<ViewStatus>('loading')
  const [calendars, setCalendars] = useState<CalendarItem[]>([])
  const [calendarsLoading, setCalendarsLoading] = useState(false)
  const [showToast, setShowToast] = useState(false)
  const [weekStart, setWeekStart] = useState(() => startOfISOMonday(new Date()))
  const [showSpinner, setShowSpinner] = useState(false)
  const [totalHours, setTotalHours] = useState(0)
  const chartRef = useRef<HTMLDivElement>(null)
  const mainContentRef = useRef<HTMLDivElement>(null)
  const currentWeekStart = useMemo(() => startOfISOMonday(new Date()), [])
  const isCurrentWeek = weekStart.getTime() >= currentWeekStart.getTime()
  const weekLabel = formatWeekLabel(weekStart)

  const fetchSpendings = useCallback(async () => {
    setLoading(true)
    setError(null)
    setErrorKind('generic')
    try {
      const startDate = formatApiDate(weekStart)
      const endDate = formatApiDate(addDays(weekStart, 6))

      const token = getIdToken()
      if (!token) {
        clearSession()
        navigate('/auth', { replace: true })
        return
      }

      const response = await fetch(`/api/calendar/spending?start=${startDate}&end=${endDate}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      if (response.status === 401) {
        clearSession()
        navigate('/auth', { replace: true })
        return
      }

      if (response.status === 424) {
        setCalendarStatus('disconnected')
        return
      }

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`)
      }
      const spendings: Spendings = await response.json()
      const chartData = toChartData(spendings)
      setData(chartData)
      setTotalHours(chartData.reduce((acc, item) => acc + item.hours, 0))
    } catch (err) {
      const offline =
        !navigator.onLine ||
        (err instanceof TypeError && err.message.toLowerCase().includes('fetch'))

      setErrorKind(offline ? 'offline' : 'generic')
      trackEvent('timespending_load_failed', {
        kind: offline ? 'offline' : 'generic',
      })
      setError(err instanceof Error ? err.message : 'Failed to fetch data')
    } finally {
      setLoading(false)
    }
  }, [navigate, weekStart])

  const loadCalendarPicker = useCallback(async () => {
    setCalendarsLoading(true)
    try {
      const list = await getCalendars()
      setCalendars(list)
      setCalendarStatus('pending_selection')
    } catch (err) {
      if (err instanceof Error && err.message.includes('401')) {
        clearSession()
        navigate('/auth', { replace: true })
        return
      }
      if (err instanceof Error && err.message.includes('424')) {
        setCalendarStatus('disconnected')
        return
      }
      setCalendars([])
      setCalendarStatus('pending_selection')
    } finally {
      setCalendarsLoading(false)
    }
  }, [navigate])

  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const forcePicker = params.get('calendar') === 'select'

    async function bootstrap() {
      try {
        const status = await getCalendarStatus()
        if (forcePicker || status === 'pending_selection') {
          await loadCalendarPicker()
          return
        }
        setCalendarStatus(status)
      } catch {
        setCalendarStatus('disconnected')
      }
    }

    bootstrap()
  }, [loadCalendarPicker])

  useEffect(() => {
    if (calendarStatus === 'connected') {
      fetchSpendings()
    }
  }, [calendarStatus, fetchSpendings])

  useEffect(() => {
    if (!loading) {
      setShowSpinner(false)
      return
    }

    const timer = window.setTimeout(() => {
      setShowSpinner(true)
    }, 200)

    return () => {
      window.clearTimeout(timer)
    }
  }, [loading])

  async function handleConnectCalendar() {
    try {
      const authURL = await getCalendarAuthURL()
      window.location.href = authURL
    } catch (err) {
      if (err instanceof Error && err.message.includes('401')) {
        clearSession()
        navigate('/auth', { replace: true })
      }
    }
  }

  async function handleSaveCalendar(calendarID: string) {
    await saveCalendarConnection(calendarID)
    setShowToast(true)
    setCalendarStatus('connected')
    window.history.replaceState({}, '', '/timespending')
    mainContentRef.current?.focus()
  }

  const handleRefresh = () => {
    if (chartRef.current) {
      chartRef.current.classList.add('refreshing')
      setTimeout(() => chartRef.current?.classList.remove('refreshing'), 600)
    }
    fetchSpendings()
  }

  const handlePreviousWeek = () => {
    setWeekStart((prev) => addDays(prev, -7))
  }

  const handleNextWeek = () => {
    if (isCurrentWeek) {
      return
    }
    setWeekStart((prev) => addDays(prev, 7))
  }

  return (
    <div className="app">
      <div className="ambient-glow ambient-glow-1" />
      <div className="ambient-glow ambient-glow-2" />
      <div className="grain-overlay" />

      <main className="main-content" ref={mainContentRef} tabIndex={-1}>
        {calendarStatus === 'loading' && (
          <section className="chart-section">
            <div className="loading-state">
              <div className="loading-spinner">
                <div className="spinner-ring" />
                <div className="spinner-ring" />
                <div className="spinner-ring" />
              </div>
              <p>Checking calendar connection...</p>
            </div>
          </section>
        )}

        {calendarStatus === 'disconnected' && (
          <CalendarConnectPrompt onConnect={handleConnectCalendar} />
        )}

        {calendarStatus === 'pending_selection' && (
          <CalendarPicker calendars={calendars} isLoading={calendarsLoading} onSave={handleSaveCalendar} />
        )}

        {calendarStatus === 'connected' && (
          <>
            <section className="week-nav-card">
              <button
                type="button"
                className="week-nav-btn"
                onClick={handlePreviousWeek}
                aria-label="Previous week"
              >
                ← Prev week
              </button>
              <p className={`week-label ${isCurrentWeek ? 'current' : ''}`}>{weekLabel}</p>
              <button
                type="button"
                className="week-nav-btn"
                onClick={handleNextWeek}
                aria-label="Next week"
                disabled={isCurrentWeek}
              >
                Next week →
              </button>
            </section>

            <section className="chart-section" ref={chartRef}>
              {showSpinner && (
                <div className="loading-state" aria-hidden="true">
                  <div className="loading-spinner">
                    <div className="spinner-ring" />
                    <div className="spinner-ring" />
                    <div className="spinner-ring" />
                  </div>
                  <p>Gathering your energy data...</p>
                </div>
              )}

              {error && (
                <div className="error-state">
                  <div className="error-icon">
                    <svg viewBox="0 0 24 24" width="48" height="48">
                      <circle cx="12" cy="12" r="10" fill="none" stroke="currentColor" strokeWidth="1.5" />
                      <path d="M12 8v4M12 16h.01" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
                    </svg>
                  </div>
                  <p className="error-message">
                    {errorKind === 'offline' ? 'You are offline' : 'Unable to load data'}
                  </p>
                  <p className="error-detail">{error}</p>
                  <button className="retry-btn" onClick={handleRefresh}>Try again</button>
                </div>
              )}

              {!loading && !error && data.length > 0 && (
                <div className="chart-container">
                  <div className="chart-header">
                    <h2>Time Distribution</h2>
                    <p className="chart-subtitle">{`${weekLabel} · ${totalHours.toFixed(1)}h total`}</p>
                  </div>

                  <div className="chart-wrapper">
                    <ResponsiveContainer width="100%" height={360}>
                      <BarChart
                        data={data}
                        margin={{ top: 20, right: 20, left: -10, bottom: 60 }}
                        barCategoryGap="20%"
                      >
                        <defs>
                          {data.map((entry, index) => (
                            <linearGradient key={`gradient-${index}`} id={`barGradient-${index}`} x1="0" y1="0" x2="0" y2="1">
                              <stop offset="0%" stopColor={getSpendingColor(entry.name, index)} stopOpacity={1} />
                              <stop offset="100%" stopColor={getSpendingColor(entry.name, index)} stopOpacity={0.6} />
                            </linearGradient>
                          ))}
                          <filter id="glow">
                            <feGaussianBlur stdDeviation="3" result="coloredBlur" />
                            <feMerge>
                              <feMergeNode in="coloredBlur" />
                              <feMergeNode in="SourceGraphic" />
                            </feMerge>
                          </filter>
                        </defs>
                        <XAxis
                          dataKey="name"
                          axisLine={false}
                          tickLine={false}
                          tick={{ fill: '#a09a90', fontSize: 13, fontWeight: 500 }}
                          tickFormatter={(value) => truncateAxisLabel(String(value))}
                          dy={10}
                          angle={-35}
                          textAnchor="end"
                        />
                        <YAxis
                          axisLine={false}
                          tickLine={false}
                          tick={{ fill: '#6b665c', fontSize: 12 }}
                          tickFormatter={(value) => `${value}h`}
                          width={45}
                        />
                        <Tooltip
                          content={<CustomTooltip />}
                          cursor={{ fill: 'rgba(232, 164, 69, 0.08)', radius: 8 }}
                        />
                        <Bar
                          dataKey="hours"
                          radius={[8, 8, 0, 0]}
                          filter="url(#glow)"
                        >
                          {data.map((_, index) => (
                            <Cell
                              key={`cell-${index}`}
                              fill={`url(#barGradient-${index})`}
                              className="bar-cell"
                            />
                          ))}
                        </Bar>
                      </BarChart>
                    </ResponsiveContainer>
                  </div>

                  <div className="category-legend">
                    {data.map((entry, index) => (
                      <div key={entry.name} className="legend-item">
                        <span
                          className="legend-dot"
                          style={{ backgroundColor: getSpendingColor(entry.name, index) }}
                        />
                        <span className="legend-name">{entry.name}</span>
                        <span className="legend-hours">{entry.hours.toFixed(1)}h</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {!loading && !error && data.length === 0 && (
                <div className="empty-state">
                  <div className="empty-icon">
                    <svg viewBox="0 0 64 64" width="64" height="64">
                      <circle cx="32" cy="32" r="28" fill="none" stroke="currentColor" strokeWidth="1.5" opacity="0.3" />
                      <path d="M20 32h24M32 20v24" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" opacity="0.5" />
                    </svg>
                  </div>
                  <p className="empty-title">No data yet</p>
                  <p className="empty-subtitle">No calendar events were found for this week.</p>
                </div>
              )}
            </section>

            <footer className="footer">
              <p>Track mindfully · Live intentionally</p>
            </footer>
          </>
        )}
      </main>

      {showToast && (
        <Toast
          message="Calendar connected!"
          subtitle="Your events are syncing…"
          onDismiss={() => setShowToast(false)}
        />
      )}
    </div>
  )
}

export default App
