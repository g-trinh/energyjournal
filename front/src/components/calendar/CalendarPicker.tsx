import { useEffect, useState } from 'react'
import type { CalendarItem } from './types'

interface CalendarPickerProps {
  calendars: CalendarItem[]
  isLoading: boolean
  onSave: (calendarID: string) => Promise<void> | void
}

export default function CalendarPicker({ calendars, isLoading, onSave }: CalendarPickerProps) {
  const [selectedID, setSelectedID] = useState('')
  const [isSaving, setIsSaving] = useState(false)

  useEffect(() => {
    if (calendars.length > 0 && !selectedID) {
      setSelectedID(calendars[0].id)
    }
  }, [calendars, selectedID])

  const isEmpty = !isLoading && calendars.length === 0
  const saveDisabled = isLoading || isSaving || isEmpty || !selectedID

  async function handleSave() {
    if (saveDisabled) {
      return
    }
    setIsSaving(true)
    try {
      await onSave(selectedID)
    } finally {
      setIsSaving(false)
    }
  }

  return (
    <section className="calendar-picker-shell">
      <div className="calendar-picker-card">
        <h2>Select a calendar</h2>
        <p className="calendar-picker-subtitle">Choose the calendar used for your time distribution.</p>

        <div className="calendar-picker-controls">
          <label htmlFor="calendar-select">Select a calendar</label>
          <select
            id="calendar-select"
            value={selectedID}
            onChange={(event) => setSelectedID(event.target.value)}
            disabled={isLoading || isSaving || isEmpty}
          >
            {isLoading && <option>Loading calendars...</option>}
            {!isLoading && calendars.map((calendar) => (
              <option key={calendar.id} value={calendar.id}>
                {calendar.name}
              </option>
            ))}
          </select>

          <button type="button" className="calendar-picker-save" onClick={handleSave} disabled={saveDisabled}>
            {isSaving ? 'Saving...' : 'Save'}
          </button>
        </div>

        {isEmpty && <p className="calendar-picker-empty">No calendars found</p>}
      </div>
    </section>
  )
}
