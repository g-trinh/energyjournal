import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { getIdToken } from '@/lib/session'
import {
  addDays,
  daysBetween,
  formatDisplayDate,
  getEnergyLevelsRange,
  type EnergyLevels,
} from '@/services/energyLevels'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { Link } from 'react-router-dom'
import '../App.css'
import '@/styles/energy-levels.css'

type Preset = '7d' | '14d' | '30d' | 'custom'

type PageStatus = 'idle' | 'loading' | 'success' | 'error' | 'empty'

const PRESET_OFFSETS: Record<Exclude<Preset, 'custom'>, number> = {
  '7d': 6,
  '14d': 13,
  '30d': 29,
}

function todayAsDateInputValue(): string {
  return new Date().toLocaleDateString('en-CA')
}

export default function EnergyLevelsPage() {
  const today = useMemo(() => todayAsDateInputValue(), [])
  const [from, setFrom] = useState<string>(() => addDays(today, -13))
  const [to, setTo] = useState<string>(() => today)
  const [preset, setPreset] = useState<Preset>('14d')
  const [levels, setLevels] = useState<EnergyLevels[]>([])
  const [status, setStatus] = useState<PageStatus>('idle')
  const [clampWarning, setClampWarning] = useState<boolean>(false)
  const fetchAbortRef = useRef<AbortController | null>(null)

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

  const isLoading = status === 'loading'

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
              <p className="energy-levels-chart-subtitle">Select a date range to view your trends.</p>
            </div>
          </CardHeader>
          <CardContent className="energy-levels-chart-body">
            {status === 'loading' && (
              <div className="loading-state">
                <div className="loading-spinner" aria-hidden="true">
                  <div className="spinner-ring" />
                  <div className="spinner-ring" />
                  <div className="spinner-ring" />
                </div>
                <p>Gathering your energy data...</p>
              </div>
            )}
            {status !== 'loading' && (
              <div className="energy-levels-chart-placeholder">
                Chart will render here.
              </div>
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
