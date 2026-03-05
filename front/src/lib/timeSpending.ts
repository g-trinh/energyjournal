import { CALENDAR_COLOR_MAP, CALENDAR_FALLBACK_COLORS } from '@/lib/calendarColors'

export type Spendings = Record<string, number>

export interface ChartData {
  name: string
  hours: number
}

export function getSpendingColor(name: string, index: number): string {
  return CALENDAR_COLOR_MAP[name] || CALENDAR_FALLBACK_COLORS[index % CALENDAR_FALLBACK_COLORS.length]
}

export function toChartData(spendings: Spendings): ChartData[] {
  return Object.entries(spendings)
    .filter(([name, hours]) => typeof name === 'string' && typeof hours === 'number' && Number.isFinite(hours))
    .map(([name, hours]) => ({ name, hours }))
    .sort((a, b) => b.hours - a.hours)
}

export function truncateAxisLabel(value: string, maxLength = 10): string {
  if (value.length <= maxLength) {
    return value
  }
  return `${value.slice(0, maxLength)}…`
}
