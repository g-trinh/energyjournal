import { describe, expect, it } from 'vitest'

import { CALENDAR_COLOR_MAP, CALENDAR_FALLBACK_COLORS } from '@/lib/calendarColors'
import { getSpendingColor, toChartData, truncateAxisLabel } from '@/lib/timeSpending'

describe('time spending helpers', () => {
  it('uses mapped calendar color when the label is known', () => {
    expect(getSpendingColor('Tomato', 0)).toBe(CALENDAR_COLOR_MAP.Tomato)
  })

  it('cycles fallback colors for unknown labels', () => {
    expect(getSpendingColor('UnknownColorA', 0)).toBe(CALENDAR_FALLBACK_COLORS[0])
    expect(getSpendingColor('UnknownColorB', 1)).toBe(CALENDAR_FALLBACK_COLORS[1])
  })

  it('returns sorted chart data and drops invalid values', () => {
    const data = toChartData({ Tomato: 2, Sage: 5, Broken: Number.NaN })
    expect(data).toEqual([
      { name: 'Sage', hours: 5 },
      { name: 'Tomato', hours: 2 },
    ])
  })

  it('truncates long axis labels', () => {
    expect(truncateAxisLabel('BlueberryLongName')).toBe('BlueberryL…')
  })
})
