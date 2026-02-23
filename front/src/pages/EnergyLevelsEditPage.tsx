import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { getIdToken } from '@/lib/session'
import {
  ENERGY_LEVELS_FORCE_REFRESH_KEY,
  clearEnergyLevelsRangeCache,
} from '@/lib/energyLevelsCache'
import {
  formatDisplayDate,
  getEnergyLevels,
  saveEnergyLevels,
  type EnergyLevels,
} from '@/services/energyLevels'
import EnergySection from '@/components/energy/EnergySection'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { createPortal } from 'react-dom'
import { useNavigate } from 'react-router-dom'
import '../App.css'
import '@/styles/energy-edit.css'

type PageStatus = 'loading' | 'idle' | 'saving'

const DEFAULT_LEVEL = 5

function todayAsDateInputValue(): string {
  return new Date().toLocaleDateString('en-CA')
}

function formatToastDate(date: string): string {
  return new Date(`${date}T00:00:00`).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

export default function EnergyLevelsEditPage() {
  const navigate = useNavigate()
  const [date, setDate] = useState<string>(todayAsDateInputValue)
  const [physical, setPhysical] = useState<number>(DEFAULT_LEVEL)
  const [mental, setMental] = useState<number>(DEFAULT_LEVEL)
  const [emotional, setEmotional] = useState<number>(DEFAULT_LEVEL)
  const [status, setStatus] = useState<PageStatus>('loading')
  const [hasExistingData, setHasExistingData] = useState<boolean>(false)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [saveError, setSaveError] = useState<string | null>(null)
  const [toastData, setToastData] = useState<EnergyLevels | null>(null)
  const fetchAbortRef = useRef<AbortController | null>(null)

  const levels = useMemo<EnergyLevels>(
    () => ({ date, physical, mental, emotional }),
    [date, physical, mental, emotional],
  )

  const loadLevels = useCallback(async (selectedDate: string) => {
    const token = getIdToken()
    if (!token) {
      setStatus('idle')
      setLoadError('Could not load data for this date. Please try again.')
      return
    }

    fetchAbortRef.current?.abort()
    const controller = new AbortController()
    fetchAbortRef.current = controller

    setStatus('loading')
    setLoadError(null)

    try {
      const result = await getEnergyLevels(selectedDate, token, controller.signal)
      if (result === null) {
        setPhysical(DEFAULT_LEVEL)
        setMental(DEFAULT_LEVEL)
        setEmotional(DEFAULT_LEVEL)
        setHasExistingData(false)
      } else {
        setPhysical(result.physical)
        setMental(result.mental)
        setEmotional(result.emotional)
        setHasExistingData(true)
      }
    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        return
      }

      setPhysical(DEFAULT_LEVEL)
      setMental(DEFAULT_LEVEL)
      setEmotional(DEFAULT_LEVEL)
      setHasExistingData(false)
      setLoadError('Could not load data for this date. Please try again.')
    } finally {
      if (!controller.signal.aborted) {
        setStatus('idle')
      }
    }
  }, [])

  useEffect(() => {
    void loadLevels(date)
    return () => {
      fetchAbortRef.current?.abort()
    }
  }, [date, loadLevels])

  useEffect(() => {
    if (!toastData) {
      return
    }

    const timeoutID = window.setTimeout(() => {
      setToastData(null)
    }, 3000)

    return () => {
      window.clearTimeout(timeoutID)
    }
  }, [toastData])

  async function handleSave() {
    const token = getIdToken()
    if (!token) {
      setSaveError('Failed to save. Please try again.')
      return
    }

    setStatus('saving')
    setSaveError(null)

    try {
      const saved = await saveEnergyLevels(levels, token)
      setPhysical(saved.physical)
      setMental(saved.mental)
      setEmotional(saved.emotional)
      setHasExistingData(true)
      setToastData(saved)
      clearEnergyLevelsRangeCache()
      window.sessionStorage.setItem(ENERGY_LEVELS_FORCE_REFRESH_KEY, '1')
      setStatus('idle')
    } catch {
      setSaveError('Failed to save. Please try again.')
      setStatus('idle')
    }
  }

  const isBusy = status === 'loading'

  function goBack() {
    navigate('/energy/levels')
  }

  return (
    <main className="app energy-edit-page">
      <div className="ambient-glow ambient-glow-1" />
      <div className="ambient-glow ambient-glow-2" />
      <div className="grain-overlay" />

      <Card className="energy-edit-card">
        <CardHeader className="energy-edit-card-header">
          <button
            type="button"
            className="energy-edit-back-btn"
            onClick={goBack}
            aria-label="Back to energy levels"
          >
            <svg width="16" height="16" viewBox="0 0 20 20" aria-hidden="true">
              <path d="M8 5l-5 5 5 5M15 10H4" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
            Back
          </button>
          <h1 className="energy-edit-title">Edit Energy Levels</h1>
          <p className="energy-edit-subtitle">Save and return to refreshed trend chart.</p>
        </CardHeader>

        <CardContent className="energy-edit-card-content">
          <hr className="energy-edit-header-divider" aria-hidden="true" />
          <div className="energy-edit-date-section">
            <label htmlFor="energy-date" className="energy-edit-date-label">
              DATE
            </label>
            <div className="energy-edit-date-input-wrap">
              <svg
                className="energy-edit-date-icon"
                width="16"
                height="16"
                viewBox="0 0 20 20"
                aria-hidden="true"
              >
                <rect x="3" y="4" width="14" height="13" rx="2" fill="none" stroke="currentColor" strokeWidth="1.5" />
                <path d="M6 2.8V6M14 2.8V6M3 8h14" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
              </svg>
              <span className="energy-edit-date-display" aria-hidden="true">
                {formatDisplayDate(date)}
              </span>
              <input
                id="energy-date"
                type="date"
                className="energy-edit-date-input"
                value={date}
                onChange={(event) => {
                  setDate(event.target.value)
                }}
                aria-label="Date"
              />
              <svg
                className="energy-edit-date-chevron"
                width="16"
                height="16"
                viewBox="0 0 20 20"
                aria-hidden="true"
              >
                <path d="M5 7l5 6 5-6" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
              </svg>
            </div>
            {loadError && <p className="energy-load-error">{loadError}</p>}
          </div>

          {hasExistingData && (
            <div className="energy-edit-status-badge" aria-live="polite">
              <svg width="16" height="16" viewBox="0 0 20 20" aria-hidden="true">
                <circle cx="10" cy="10" r="8" fill="#8fa58b" />
                <path d="M6 10.2l2.4 2.4L14.2 7" fill="none" stroke="#f5f0e8" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
              <span className="energy-status-text-desktop">
                Existing levels loaded for this date — adjust and save to update
              </span>
              <span className="energy-status-text-mobile">
                Existing levels loaded - tap to adjust
              </span>
            </div>
          )}

          {status === 'loading' && <p className="energy-loading">Loading energy levels...</p>}

          <EnergySection
            label="Physical"
            color="#c4826d"
            value={physical}
            onChange={setPhysical}
            disabled={isBusy}
          />
          <EnergySection
            label="Mental"
            color="#7eb8b3"
            value={mental}
            onChange={setMental}
            disabled={isBusy}
          />
          <EnergySection
            label="Emotional"
            color="#8fa58b"
            value={emotional}
            onChange={setEmotional}
            showDivider={false}
            disabled={isBusy}
          />

          <div className="energy-edit-actions">
            <button
              type="button"
              className="energy-save-btn"
              disabled={status === 'loading' || status === 'saving'}
              aria-disabled={status === 'loading' || status === 'saving'}
              onClick={() => {
                void handleSave()
              }}
            >
              {status === 'saving' ? 'Saving…' : 'Save Energy Levels'}
            </button>
            {saveError && <p className="energy-save-error">{saveError}</p>}
            <button
              className="energy-cancel-link"
              type="button"
              onClick={goBack}
            >
              Cancel
            </button>
          </div>
        </CardContent>
      </Card>

      {toastData &&
        createPortal(
          <div className="energy-toast" role="status" aria-live="polite">
            <svg width="16" height="16" viewBox="0 0 20 20" aria-hidden="true">
              <circle cx="10" cy="10" r="8" fill="#8fa58b" />
              <path d="M6 10.2l2.4 2.4L14.2 7" fill="none" stroke="#f5f0e8" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
            <div className="energy-toast-copy">
              <p className="energy-toast-title">Energy levels saved</p>
              <p className="energy-toast-subtitle">
                {formatToastDate(toastData.date)}
                {' · Physical '}
                {toastData.physical}
                {' · Mental '}
                {toastData.mental}
                {' · Emotional '}
                {toastData.emotional}
              </p>
            </div>
            <button
              type="button"
              className="energy-toast-close"
              onClick={() => {
                setToastData(null)
              }}
              aria-label="Dismiss save confirmation"
            >
              ×
            </button>
          </div>,
          document.body,
        )}
    </main>
  )
}
