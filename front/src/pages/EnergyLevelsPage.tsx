import { Card, CardContent, CardHeader } from '@/components/ui/card'
import EnergyTooltip from '@/components/energy/EnergyTooltip'
import { ENERGY_COLORS, ENERGY_LABELS, type EnergyDimension } from '@/lib/energyColors'
import { getIdToken } from '@/lib/session'
import {
  addDays,
  daysBetween,
  formatDisplayDate,
  getDaysBetween,
  getEnergyLevelsRange,
  type EnergyLevels,
} from '@/services/energyLevels'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { Link } from 'react-router-dom'
import {
  CartesianGrid,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import '../App.css'
import '@/styles/energy-levels.css'

type Preset = '7d' | '14d' | '30d' | 'custom'

type PageStatus = 'idle' | 'loading' | 'success' | 'error' | 'empty'

interface ChartPoint {
  date: string
  physical: number | null
  mental: number | null
  emotional: number | null
}

const PRESET_OFFSETS: Record<Exclude<Preset, 'custom'>, number> = {
  '7d': 6,
  '14d': 13,
  '30d': 29,
}

function todayAsDateInputValue(): string {
  return new Date().toLocaleDateString('en-CA')
}

function formatShortDate(date: string): string {
  return new Date(`${date}T00:00:00`).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
  })
}

function formatRangeLabel(from: string, to: string): string {
  const fromDate = new Date(`${from}T00:00:00`)
  const toDate = new Date(`${to}T00:00:00`)
  const fromLabel = formatShortDate(from)
  const toLabel = toDate.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })

  if (fromDate.getFullYear() !== toDate.getFullYear()) {
    const fromWithYear = fromDate.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    })
    return `${fromWithYear} - ${toLabel}`
  }

  return `${fromLabel} - ${toLabel}`
}

export function buildChartData(levels: EnergyLevels[], from: string, to: string): ChartPoint[] {
  if (!from || !to) {
    return []
  }
  const map = new Map(levels.map((level) => [level.date, level]))
  return getDaysBetween(from, to).map((date) => ({
    date,
    physical: map.get(date)?.physical ?? null,
    mental: map.get(date)?.mental ?? null,
    emotional: map.get(date)?.emotional ?? null,
  }))
}

export default function EnergyLevelsPage() {
  const today = useMemo(() => todayAsDateInputValue(), [])
  const [from, setFrom] = useState<string>(() => addDays(today, -13))
  const [to, setTo] = useState<string>(() => today)
  const [preset, setPreset] = useState<Preset>('14d')
  const [levels, setLevels] = useState<EnergyLevels[]>([])
  const [status, setStatus] = useState<PageStatus>('idle')
  const [clampWarning, setClampWarning] = useState<boolean>(false)
  const [hiddenLines, setHiddenLines] = useState<Set<EnergyDimension>>(new Set())
  const [isMobile, setIsMobile] = useState<boolean>(() => {
    if (typeof window === 'undefined' || !('matchMedia' in window)) {
      return false
    }
    return window.matchMedia('(max-width: 768px)').matches
  })
  const fetchAbortRef = useRef<AbortController | null>(null)

  useEffect(() => {
    if (typeof window === 'undefined' || !('matchMedia' in window)) {
      return undefined
    }

    const media = window.matchMedia('(max-width: 768px)')
    const handler = (event: MediaQueryListEvent) => setIsMobile(event.matches)
    if ('addEventListener' in media) {
      media.addEventListener('change', handler)
    } else {
      media.addListener(handler)
    }

    return () => {
      if ('removeEventListener' in media) {
        media.removeEventListener('change', handler)
      } else {
        media.removeListener(handler)
      }
    }
  }, [])

  const loadRange = useCallback(async (nextFrom: string, nextTo: string) => {
    const token = getIdToken()
    if (!token) {
      setStatus('error')
      return
    }

    fetchAbortRef.current?.abort()
    const controller = new AbortController()
    fetchAbortRef.current = controller

    setStatus('loading')

    try {
      const result = await getEnergyLevelsRange(
        nextFrom,
        nextTo,
        token,
        controller.signal,
      )
      if (result.length === 0) {
        setLevels([])
        setStatus('empty')
      } else {
        setLevels(result)
        setStatus('success')
      }
    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        return
      }

      setStatus('error')
    }
  }, [])

  useEffect(() => {
    void loadRange(from, to)
    return () => {
      fetchAbortRef.current?.abort()
    }
  }, [from, to, loadRange])

  const chartData = useMemo(
    () => buildChartData(levels, from, to),
    [levels, from, to],
  )
  const dayCount = from && to ? Math.max(daysBetween(from, to) + 1, 1) : 0
  const chartSubtitle = from && to
    ? `${formatRangeLabel(from, to)} · ${dayCount} days`
    : ''
  const tickInterval = useMemo(() => {
    if (chartData.length === 0) {
      return 0
    }
    const targetTicks = isMobile ? 4 : 7
    return Math.max(Math.ceil(chartData.length / targetTicks) - 1, 0)
  }, [chartData.length, isMobile])

  function applyPreset(value: Exclude<Preset, 'custom'>) {
    const nextTo = todayAsDateInputValue()
    const nextFrom = addDays(nextTo, -PRESET_OFFSETS[value])
    setFrom(nextFrom)
    setTo(nextTo)
    setPreset(value)
    setClampWarning(false)
  }

  function handleDateChange(field: 'from' | 'to', value: string) {
    let nextFrom = field === 'from' ? value : from
    let nextTo = field === 'to' ? value : to

    if (nextFrom && nextTo && daysBetween(nextFrom, nextTo) < 0) {
      nextTo = nextFrom
    }

    const diff = nextFrom && nextTo ? daysBetween(nextFrom, nextTo) : 0
    if (diff > 30) {
      if (field === 'from') {
        nextTo = addDays(nextFrom, 30)
      } else {
        nextFrom = addDays(nextTo, -30)
      }
      setClampWarning(true)
    } else {
      setClampWarning(false)
    }

    setFrom(nextFrom)
    setTo(nextTo)
    setPreset('custom')
  }

  function toggleLine(dimension: EnergyDimension) {
    setHiddenLines((prev) => {
      const next = new Set(prev)
      if (next.has(dimension)) {
        next.delete(dimension)
      } else {
        next.add(dimension)
      }
      return next
    })
  }

  const isLoading = status === 'loading'
  const showChart = status === 'success'
  const showEmpty = status === 'empty'
  const showError = status === 'error'

  return (
    <main className="app energy-levels-page">
      <div className="ambient-glow ambient-glow-1" />
      <div className="ambient-glow ambient-glow-2" />
      <div className="grain-overlay" />

      <section className="energy-levels-content">
        <header className="energy-levels-header">
          <div>
            <h1 className="energy-levels-title">Energy Levels</h1>
            <p className="energy-levels-subtitle">
              Track how your physical, mental, and emotional energy evolves over time.
            </p>
          </div>
          <Link to="/energy/levels/edit" className="landing-cta-primary energy-levels-track-btn-desktop">
            Track Today
          </Link>
        </header>

        <Card className="energy-levels-date-range-card">
          <CardHeader className="energy-levels-date-header">
            <div className="energy-levels-presets">
              {(['7d', '14d', '30d'] as const).map((value) => (
                <button
                  key={value}
                  type="button"
                  className={
                    preset === value
                      ? 'range-preset-btn range-preset-btn-active'
                      : 'range-preset-btn'
                  }
                  aria-pressed={preset === value}
                  onClick={() => applyPreset(value)}
                >
                  {value}
                </button>
              ))}
            </div>
          </CardHeader>
          <CardContent className="energy-levels-date-body">
            <div className="energy-levels-date-inputs">
              <div className="energy-levels-date-field">
                <label htmlFor="energy-range-from">FROM</label>
                <div className="energy-levels-date-input-wrap">
                  <svg
                    className="energy-levels-date-icon"
                    width="16"
                    height="16"
                    viewBox="0 0 20 20"
                    aria-hidden="true"
                  >
                    <rect x="3" y="4" width="14" height="13" rx="2" fill="none" stroke="currentColor" strokeWidth="1.5" />
                    <path d="M6 2.8V6M14 2.8V6M3 8h14" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
                  </svg>
                  <span className="energy-levels-date-display" aria-hidden="true">
                    {formatDisplayDate(from)}
                  </span>
                  <input
                    id="energy-range-from"
                    type="date"
                    className="energy-levels-date-input"
                    value={from}
                    onChange={(event) => handleDateChange('from', event.target.value)}
                  />
                  <svg
                    className="energy-levels-date-chevron"
                    width="16"
                    height="16"
                    viewBox="0 0 20 20"
                    aria-hidden="true"
                  >
                    <path d="M5 7l5 6 5-6" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
                  </svg>
                </div>
              </div>

              <div className="energy-levels-date-separator" aria-hidden="true">
                &rarr;
              </div>

              <div className="energy-levels-date-field">
                <label htmlFor="energy-range-to">TO</label>
                <div className="energy-levels-date-input-wrap">
                  <svg
                    className="energy-levels-date-icon"
                    width="16"
                    height="16"
                    viewBox="0 0 20 20"
                    aria-hidden="true"
                  >
                    <rect x="3" y="4" width="14" height="13" rx="2" fill="none" stroke="currentColor" strokeWidth="1.5" />
                    <path d="M6 2.8V6M14 2.8V6M3 8h14" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
                  </svg>
                  <span className="energy-levels-date-display" aria-hidden="true">
                    {formatDisplayDate(to)}
                  </span>
                  <input
                    id="energy-range-to"
                    type="date"
                    className="energy-levels-date-input"
                    value={to}
                    onChange={(event) => handleDateChange('to', event.target.value)}
                  />
                  <svg
                    className="energy-levels-date-chevron"
                    width="16"
                    height="16"
                    viewBox="0 0 20 20"
                    aria-hidden="true"
                  >
                    <path d="M5 7l5 6 5-6" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
                  </svg>
                </div>
              </div>

              <button
                type="button"
                className={
                  isLoading
                    ? 'energy-levels-refresh-btn energy-levels-refresh-loading'
                    : 'energy-levels-refresh-btn'
                }
                onClick={() => void loadRange(from, to)}
                aria-label="Refresh range"
              >
                <svg
                  width="16"
                  height="16"
                  viewBox="0 0 20 20"
                  aria-hidden="true"
                >
                  <path
                    d="M15 6.5a5 5 0 10.8 5.7"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="1.6"
                    strokeLinecap="round"
                  />
                  <path d="M15 3v3.8h-3.8" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" />
                </svg>
                <span>Refresh</span>
              </button>
            </div>
            {clampWarning && (
              <p className="clamp-warning" role="status" aria-live="polite">
                Range limited to 30 days
              </p>
            )}
          </CardContent>
        </Card>

        <Card className="energy-levels-chart-card" aria-busy={isLoading}>
          <CardHeader className="energy-levels-chart-header">
            <div>
              <h2>Your Energy Trends</h2>
              <p className="energy-levels-chart-subtitle">
                {chartSubtitle || 'Select a date range to view your trends.'}
              </p>
            </div>
          </CardHeader>
          <CardContent className="energy-levels-chart-body">
            {(isLoading || status === 'idle') && (
              <div className="loading-state">
                <div className="loading-spinner" aria-hidden="true">
                  <div className="spinner-ring" />
                  <div className="spinner-ring" />
                  <div className="spinner-ring" />
                </div>
                <p>Gathering your energy data...</p>
              </div>
            )}

            {showError && (
              <div className="error-state" role="alert">
                <svg className="error-icon" width="32" height="32" viewBox="0 0 20 20" aria-hidden="true">
                  <circle cx="10" cy="10" r="8" fill="none" stroke="currentColor" strokeWidth="1.6" />
                  <path d="M10 6v5M10 14v.5" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" />
                </svg>
                <p className="error-message">Unable to load data</p>
                <p className="error-detail">Please check your connection and try again.</p>
                <button type="button" className="retry-btn" onClick={() => void loadRange(from, to)}>
                  Try again
                </button>
              </div>
            )}

            {showEmpty && (
              <div className="empty-state">
                <svg className="empty-icon" width="36" height="36" viewBox="0 0 20 20" aria-hidden="true">
                  <path d="M4 10h12" fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
                  <circle cx="10" cy="10" r="7" fill="none" stroke="currentColor" strokeWidth="1.4" />
                </svg>
                <p className="empty-title">No entries yet</p>
                <p className="empty-subtitle">Start tracking your energy to see your trends here.</p>
                <Link to="/energy/levels/edit" className="energy-levels-empty-cta">
                  Track Today
                </Link>
              </div>
            )}

            {showChart && (
              <>
                <div className="energy-levels-chart-wrapper">
                  <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={chartData} margin={{ top: 8, right: 16, left: -12, bottom: 0 }}>
                      <CartesianGrid vertical={false} stroke="#a09a90" strokeOpacity={0.06} />
                      <XAxis
                        dataKey="date"
                        tickFormatter={formatShortDate}
                        interval={tickInterval}
                        tick={{ fill: '#a09a90', fontSize: 11 }}
                        axisLine={false}
                        tickLine={false}
                      />
                      <YAxis
                        domain={[0, 10]}
                        ticks={[0, 2, 4, 6, 8, 10]}
                        tick={{ fill: '#a09a90', fontSize: 11 }}
                        axisLine={false}
                        tickLine={false}
                      />
                      <Tooltip content={<EnergyTooltip />} />
                      <Line
                        dataKey="physical"
                        stroke={ENERGY_COLORS.physical}
                        strokeWidth={2.5}
                        connectNulls={false}
                        hide={hiddenLines.has('physical')}
                        dot={false}
                        activeDot={{ r: 5, stroke: '#26231e', strokeWidth: 2 }}
                      />
                      <Line
                        dataKey="mental"
                        stroke={ENERGY_COLORS.mental}
                        strokeWidth={2.5}
                        connectNulls={false}
                        hide={hiddenLines.has('mental')}
                        dot={false}
                        activeDot={{ r: 5 }}
                      />
                      <Line
                        dataKey="emotional"
                        stroke={ENERGY_COLORS.emotional}
                        strokeWidth={2.5}
                        connectNulls={false}
                        hide={hiddenLines.has('emotional')}
                        dot={false}
                        activeDot={{ r: 5 }}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                </div>
                <div className="energy-levels-legend" role="group" aria-label="Toggle energy dimensions">
                  {(Object.keys(ENERGY_COLORS) as EnergyDimension[]).map((dimension) => {
                    const isHidden = hiddenLines.has(dimension)
                    return (
                      <button
                        key={dimension}
                        type="button"
                        className="energy-levels-legend-btn"
                        aria-pressed={isHidden}
                        aria-label={`Toggle ${ENERGY_LABELS[dimension]} line`}
                        onClick={() => toggleLine(dimension)}
                        style={{ opacity: isHidden ? 0.3 : 1 }}
                      >
                        <span
                          className="energy-levels-legend-dot"
                          style={{ backgroundColor: ENERGY_COLORS[dimension] }}
                          aria-hidden="true"
                        />
                        {ENERGY_LABELS[dimension]}
                      </button>
                    )
                  })}
                </div>
              </>
            )}
          </CardContent>
        </Card>
      </section>

      <div className="energy-levels-sticky-cta">
        <Link to="/energy/levels/edit" className="energy-levels-track-btn">
          Track Today
        </Link>
      </div>
    </main>
  )
}
