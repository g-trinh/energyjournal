import { useCallback, useEffect, useRef, useState } from 'react'
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from 'recharts'
import { useNavigate } from 'react-router-dom'
import { trackEvent } from '@/lib/analytics'
import { clearSession, getIdToken } from '@/lib/session'
import './App.css'

type Spendings = Record<string, number>

interface ChartData {
  name: string
  hours: number
}

type FetchErrorKind = 'offline' | 'generic'

// Warm, organic color palette
const CATEGORY_COLORS: Record<string, string> = {
  'Travail': '#e8a445',      // Amber gold
  'Perso': '#8fa58b',        // Sage green
  'Routine': '#c4826d',      // Terracotta
  'Repas': '#7eb8b3',        // Seafoam
  'Sport': '#d4a574',        // Warm sand
  'Sommeil': '#9b8aa6',      // Dusty lavender
}

const DEFAULT_COLORS = ['#e8a445', '#8fa58b', '#c4826d', '#7eb8b3', '#d4a574', '#9b8aa6', '#c9a87c', '#a8b5a0']

function getColorForCategory(name: string, index: number): string {
  return CATEGORY_COLORS[name] || DEFAULT_COLORS[index % DEFAULT_COLORS.length]
}

function formatDateForDisplay(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' })
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
  const [startDate, setStartDate] = useState(() => {
    const today = new Date()
    const dayOfWeek = today.getDay()
    // Get last Sunday (end of last complete week)
    const lastSunday = new Date(today)
    lastSunday.setDate(today.getDate() - (dayOfWeek === 0 ? 7 : dayOfWeek))
    // Get the Monday before that (start of last complete week)
    const lastMonday = new Date(lastSunday)
    lastMonday.setDate(lastSunday.getDate() - 6)
    return lastMonday.toISOString().split('T')[0]
  })
  const [endDate, setEndDate] = useState(() => {
    const today = new Date()
    const dayOfWeek = today.getDay()
    // Get last Sunday (end of last complete week)
    const lastSunday = new Date(today)
    lastSunday.setDate(today.getDate() - (dayOfWeek === 0 ? 7 : dayOfWeek))
    return lastSunday.toISOString().split('T')[0]
  })
  const [totalHours, setTotalHours] = useState(0)
  const chartRef = useRef<HTMLDivElement>(null)

  const fetchSpendings = useCallback(async () => {
    setLoading(true)
    setError(null)
    setErrorKind('generic')
    try {
      const dateRegex = /^\d{4}-\d{2}-\d{2}$/
      if (!dateRegex.test(startDate) || !dateRegex.test(endDate)) {
        setError('Please use valid start and end dates.')
        return
      }

      if (startDate > endDate) {
        setError('Start date cannot be after end date.')
        return
      }

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

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`)
      }
      const spendings: Spendings = await response.json()
      const chartData = Object.entries(spendings)
        .filter(([name, hours]) => typeof name === 'string' && typeof hours === 'number' && Number.isFinite(hours))
        .map(([name, hours]) => ({ name, hours }))
        .sort((a, b) => b.hours - a.hours)
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
  }, [endDate, navigate, startDate])

  useEffect(() => {
    fetchSpendings()
  }, [fetchSpendings])

  const handleRefresh = () => {
    if (chartRef.current) {
      chartRef.current.classList.add('refreshing')
      setTimeout(() => chartRef.current?.classList.remove('refreshing'), 600)
    }
    fetchSpendings()
  }

  return (
    <div className="app">
      {/* Ambient background elements */}
      <div className="ambient-glow ambient-glow-1" />
      <div className="ambient-glow ambient-glow-2" />
      <div className="grain-overlay" />

      <main className="main-content">
        {/* Header */}
        <header className="header">
          <div className="logo-mark">
            <svg viewBox="0 0 32 32" className="logo-icon">
              <circle cx="16" cy="16" r="12" fill="none" stroke="currentColor" strokeWidth="1.5" opacity="0.3" />
              <circle cx="16" cy="16" r="8" fill="none" stroke="currentColor" strokeWidth="1.5" opacity="0.5" />
              <circle cx="16" cy="16" r="4" fill="currentColor" />
            </svg>
          </div>
          <div className="header-text">
            <h1>Energy Journal</h1>
            <p className="tagline">Where does your time flow?</p>
          </div>
        </header>

        {/* Date Range Card */}
        <section className="date-range-card">
          <div className="date-range-inner">
            <div className="date-field">
              <label htmlFor="start-date">From</label>
              <div className="date-input-wrapper">
                <input
                  id="start-date"
                  type="date"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                />
                <span className="date-display">{formatDateForDisplay(startDate)}</span>
              </div>
            </div>

            <div className="date-separator">
              <svg viewBox="0 0 24 24" width="20" height="20">
                <path d="M5 12h14M13 6l6 6-6 6" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>
            </div>

            <div className="date-field">
              <label htmlFor="end-date">To</label>
              <div className="date-input-wrapper">
                <input
                  id="end-date"
                  type="date"
                  value={endDate}
                  onChange={(e) => setEndDate(e.target.value)}
                />
                <span className="date-display">{formatDateForDisplay(endDate)}</span>
              </div>
            </div>

            <button className="refresh-btn" onClick={handleRefresh} disabled={loading}>
              <svg viewBox="0 0 24 24" width="18" height="18" className={loading ? 'spinning' : ''}>
                <path d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                <path d="M21 3v5h-5" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>
              <span>Refresh</span>
            </button>
          </div>
        </section>

        {/* Stats Summary */}
        {!loading && !error && data.length > 0 && (
          <section className="stats-summary">
            <div className="stat-card stat-total">
              <span className="stat-value">{totalHours.toFixed(1)}</span>
              <span className="stat-label">hours tracked</span>
            </div>
            <div className="stat-card stat-categories">
              <span className="stat-value">{data.length}</span>
              <span className="stat-label">categories</span>
            </div>
            <div className="stat-card stat-top">
              <span className="stat-value">{data[0]?.name || '—'}</span>
              <span className="stat-label">top activity</span>
            </div>
          </section>
        )}

        {/* Chart Section */}
        <section className="chart-section" ref={chartRef}>
          {loading && (
            <div className="loading-state">
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
                  <circle cx="12" cy="12" r="10" fill="none" stroke="currentColor" strokeWidth="1.5"/>
                  <path d="M12 8v4M12 16h.01" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
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
                <p className="chart-subtitle">Hours per category</p>
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
                          <stop offset="0%" stopColor={getColorForCategory(entry.name, index)} stopOpacity={1}/>
                          <stop offset="100%" stopColor={getColorForCategory(entry.name, index)} stopOpacity={0.6}/>
                        </linearGradient>
                      ))}
                      <filter id="glow">
                        <feGaussianBlur stdDeviation="3" result="coloredBlur"/>
                        <feMerge>
                          <feMergeNode in="coloredBlur"/>
                          <feMergeNode in="SourceGraphic"/>
                        </feMerge>
                      </filter>
                    </defs>
                    <XAxis
                      dataKey="name"
                      axisLine={false}
                      tickLine={false}
                      tick={{ fill: '#a09a90', fontSize: 13, fontWeight: 500 }}
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

              {/* Category Legend */}
              <div className="category-legend">
                {data.map((entry, index) => (
                  <div key={entry.name} className="legend-item">
                    <span
                      className="legend-dot"
                      style={{ backgroundColor: getColorForCategory(entry.name, index) }}
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
                  <circle cx="32" cy="32" r="28" fill="none" stroke="currentColor" strokeWidth="1.5" opacity="0.3"/>
                  <path d="M20 32h24M32 20v24" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" opacity="0.5"/>
                </svg>
              </div>
              <p className="empty-title">No data yet</p>
              <p className="empty-subtitle">Select a date range and refresh to see your energy flow</p>
            </div>
          )}
        </section>

        {/* Footer */}
        <footer className="footer">
          <p>Track mindfully · Live intentionally</p>
        </footer>
      </main>
    </div>
  )
}

export default App
