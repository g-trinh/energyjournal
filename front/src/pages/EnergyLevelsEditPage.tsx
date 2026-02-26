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
import ContextSelect from '@/components/energy/ContextSelect'
import StepIndicator from '@/components/energy/StepIndicator'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { createPortal } from 'react-dom'
import { useNavigate } from 'react-router-dom'
import '../App.css'
import '@/styles/energy-edit.css'
import '@/styles/energy-context.css'

type PageStatus = 'loading' | 'idle' | 'saving'
type FormStep = 1 | 2
type ContextSelectID =
  | 'physicalActivity'
  | 'nutrition'
  | 'socialInteractions'
  | 'timeOutdoors'

const DEFAULT_LEVEL = 5
const DEFAULT_CONTEXT_LEVEL = 3

const PHYSICAL_ACTIVITY_OPTIONS = [
  { value: '', label: '—' },
  { value: 'none', label: 'None' },
  { value: 'light', label: 'Light' },
  { value: 'moderate', label: 'Moderate' },
  { value: 'intense', label: 'Intense' },
]

const NUTRITION_OPTIONS = [
  { value: '', label: '—' },
  { value: 'poor', label: 'Poor quality' },
  { value: 'average', label: 'Average quality' },
  { value: 'good', label: 'Good quality' },
  { value: 'excellent', label: 'Excellent quality' },
]

const SOCIAL_INTERACTIONS_OPTIONS = [
  { value: '', label: '—' },
  { value: 'negative', label: 'Negative' },
  { value: 'neutral', label: 'Neutral' },
  { value: 'positive', label: 'Positive' },
]

const TIME_OUTDOORS_OPTIONS = [
  { value: '', label: '—' },
  { value: 'none', label: 'None' },
  { value: 'under_30min', label: 'Under 30 min' },
  { value: '30min_1hr', label: '30 min-1 hr' },
  { value: 'over_1hr', label: 'Over 1 hr' },
]

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
  const [step, setStep] = useState<FormStep>(1)
  const [date, setDate] = useState<string>(todayAsDateInputValue)
  const [physical, setPhysical] = useState<number>(DEFAULT_LEVEL)
  const [mental, setMental] = useState<number>(DEFAULT_LEVEL)
  const [emotional, setEmotional] = useState<number>(DEFAULT_LEVEL)
  const [sleepQuality, setSleepQuality] = useState<number>(DEFAULT_CONTEXT_LEVEL)
  const [stressLevel, setStressLevel] = useState<number>(DEFAULT_CONTEXT_LEVEL)
  const [physicalActivity, setPhysicalActivity] = useState<string>('')
  const [nutrition, setNutrition] = useState<string>('')
  const [socialInteractions, setSocialInteractions] = useState<string>('')
  const [timeOutdoors, setTimeOutdoors] = useState<string>('')
  const [notes, setNotes] = useState<string>('')
  const [status, setStatus] = useState<PageStatus>('loading')
  const [openContextSelect, setOpenContextSelect] = useState<ContextSelectID | null>(
    null,
  )
  const [hasExistingData, setHasExistingData] = useState<boolean>(false)
  const [loadError, setLoadError] = useState<string | null>(null)
  const [saveError, setSaveError] = useState<string | null>(null)
  const [toastData, setToastData] = useState<EnergyLevels | null>(null)
  const fetchAbortRef = useRef<AbortController | null>(null)

  const levels = useMemo<EnergyLevels>(
    () => ({
      date,
      physical,
      mental,
      emotional,
      sleepQuality,
      stressLevel,
      physicalActivity: physicalActivity || undefined,
      nutrition: nutrition || undefined,
      socialInteractions: socialInteractions || undefined,
      timeOutdoors: timeOutdoors || undefined,
      notes: notes || undefined,
    }),
    [
      date,
      emotional,
      mental,
      notes,
      nutrition,
      physical,
      physicalActivity,
      sleepQuality,
      socialInteractions,
      stressLevel,
      timeOutdoors,
    ],
  )

  const resetContextFields = useCallback(() => {
    setSleepQuality(DEFAULT_CONTEXT_LEVEL)
    setStressLevel(DEFAULT_CONTEXT_LEVEL)
    setPhysicalActivity('')
    setNutrition('')
    setSocialInteractions('')
    setTimeOutdoors('')
    setNotes('')
  }, [])

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
        resetContextFields()
        setHasExistingData(false)
      } else {
        setPhysical(result.physical)
        setMental(result.mental)
        setEmotional(result.emotional)
        setSleepQuality(result.sleepQuality ?? DEFAULT_CONTEXT_LEVEL)
        setStressLevel(result.stressLevel ?? DEFAULT_CONTEXT_LEVEL)
        setPhysicalActivity(result.physicalActivity ?? '')
        setNutrition(result.nutrition ?? '')
        setSocialInteractions(result.socialInteractions ?? '')
        setTimeOutdoors(result.timeOutdoors ?? '')
        setNotes(result.notes ?? '')
        setHasExistingData(true)
      }
    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') {
        return
      }

      setPhysical(DEFAULT_LEVEL)
      setMental(DEFAULT_LEVEL)
      setEmotional(DEFAULT_LEVEL)
      resetContextFields()
      setHasExistingData(false)
      setLoadError('Could not load data for this date. Please try again.')
    } finally {
      if (!controller.signal.aborted) {
        setStatus('idle')
      }
    }
  }, [resetContextFields])

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
      setSleepQuality(saved.sleepQuality ?? DEFAULT_CONTEXT_LEVEL)
      setStressLevel(saved.stressLevel ?? DEFAULT_CONTEXT_LEVEL)
      setPhysicalActivity(saved.physicalActivity ?? '')
      setNutrition(saved.nutrition ?? '')
      setSocialInteractions(saved.socialInteractions ?? '')
      setTimeOutdoors(saved.timeOutdoors ?? '')
      setNotes(saved.notes ?? '')
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
  const isContextDisabled = status !== 'idle'

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
          <StepIndicator
            currentStep={step}
            onStepClick={(targetStep) => {
              if (targetStep === 1) {
                setStep(1)
              }
            }}
          />
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
                  setStep(1)
                  resetContextFields()
                  setOpenContextSelect(null)
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

          {step === 1 ? (
            <>
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
                  disabled={status === 'loading'}
                  aria-disabled={status === 'loading'}
                  onClick={() => {
                    setStep(2)
                  }}
                >
                  Next Step →
                </button>
                <button
                  className="energy-cancel-link"
                  type="button"
                  onClick={goBack}
                >
                  Cancel
                </button>
              </div>
            </>
          ) : (
            <div className="energy-edit-step-two">
              <h2 className="energy-context-title">Daily Context</h2>
              <p className="energy-context-subtitle">
                Help us understand what shaped your energy today
              </p>
              <hr className="energy-context-divider" aria-hidden="true" />

              <EnergySection
                label="Sleep Quality"
                color="#7eb8b3"
                value={sleepQuality}
                onChange={setSleepQuality}
                min={1}
                max={5}
                disabled={isContextDisabled}
              />
              <EnergySection
                label="Stress Level"
                color="#c4826d"
                value={stressLevel}
                onChange={setStressLevel}
                min={1}
                max={5}
                showDivider={false}
                disabled={isContextDisabled}
              />

              <hr className="energy-context-divider energy-context-divider-loose" aria-hidden="true" />

              <div className="energy-context-select-grid">
                <ContextSelect
                  label="Physical Activity"
                  options={PHYSICAL_ACTIVITY_OPTIONS}
                  value={physicalActivity}
                  onChange={setPhysicalActivity}
                  disabled={isContextDisabled}
                  isOpen={openContextSelect === 'physicalActivity'}
                  onOpenChange={(open) =>
                    setOpenContextSelect(open ? 'physicalActivity' : null)
                  }
                />
                <ContextSelect
                  label="Nutrition"
                  options={NUTRITION_OPTIONS}
                  value={nutrition}
                  onChange={setNutrition}
                  disabled={isContextDisabled}
                  isOpen={openContextSelect === 'nutrition'}
                  onOpenChange={(open) =>
                    setOpenContextSelect(open ? 'nutrition' : null)
                  }
                />
                <ContextSelect
                  label="Social Interactions"
                  options={SOCIAL_INTERACTIONS_OPTIONS}
                  value={socialInteractions}
                  onChange={setSocialInteractions}
                  disabled={isContextDisabled}
                  isOpen={openContextSelect === 'socialInteractions'}
                  onOpenChange={(open) =>
                    setOpenContextSelect(open ? 'socialInteractions' : null)
                  }
                />
                <ContextSelect
                  label="Time Outdoors"
                  options={TIME_OUTDOORS_OPTIONS}
                  value={timeOutdoors}
                  onChange={setTimeOutdoors}
                  disabled={isContextDisabled}
                  isOpen={openContextSelect === 'timeOutdoors'}
                  onOpenChange={(open) =>
                    setOpenContextSelect(open ? 'timeOutdoors' : null)
                  }
                />
              </div>

              <hr className="energy-context-divider" aria-hidden="true" />

              <div className="energy-context-notes">
                <label htmlFor="daily-reflection" className="energy-context-label">
                  Daily Reflection
                </label>
                <textarea
                  id="daily-reflection"
                  className="energy-context-textarea"
                  placeholder="Write your thoughts about today…"
                  value={notes}
                  onChange={(event) => setNotes(event.target.value)}
                  disabled={isContextDisabled}
                  rows={5}
                />
              </div>

              <div className="energy-context-actions">
                <button
                  type="button"
                  className="energy-context-previous-btn"
                  onClick={() => setStep(1)}
                  disabled={status === 'saving'}
                >
                  ← Previous Step
                </button>
                <button
                  type="button"
                  className="energy-save-btn energy-context-save-btn"
                  disabled={status === 'loading' || status === 'saving'}
                  aria-disabled={status === 'loading' || status === 'saving'}
                  onClick={() => {
                    void handleSave()
                  }}
                >
                  {status === 'saving' ? 'Saving…' : 'Save Entry'}
                </button>
              </div>
              {saveError && <p className="energy-save-error">{saveError}</p>}
            </div>
          )}
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
