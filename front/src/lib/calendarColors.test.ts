import { CALENDAR_COLOR_MAP, CALENDAR_FALLBACK_COLORS } from './calendarColors'

describe('calendar color mapping', () => {
  it('contains all 12 Google label names', () => {
    expect(CALENDAR_COLOR_MAP.Default).toBeDefined()
    expect(CALENDAR_COLOR_MAP.Tomato).toBeDefined()
    expect(CALENDAR_COLOR_MAP.Flamingo).toBeDefined()
    expect(CALENDAR_COLOR_MAP.Tangerine).toBeDefined()
    expect(CALENDAR_COLOR_MAP.Banana).toBeDefined()
    expect(CALENDAR_COLOR_MAP.Sage).toBeDefined()
    expect(CALENDAR_COLOR_MAP.Basil).toBeDefined()
    expect(CALENDAR_COLOR_MAP.Peacock).toBeDefined()
    expect(CALENDAR_COLOR_MAP.Blueberry).toBeDefined()
    expect(CALENDAR_COLOR_MAP.Lavender).toBeDefined()
    expect(CALENDAR_COLOR_MAP.Grape).toBeDefined()
    expect(CALENDAR_COLOR_MAP.Graphite).toBeDefined()
  })

  it('provides fallback colors for unknown labels', () => {
    expect(CALENDAR_FALLBACK_COLORS.length).toBeGreaterThan(0)
  })
})
