const API_BASE = import.meta.env.VITE_API_BASE_URL ?? '/api'

export interface EnergyLevels {
  date: string
  physical: number
  mental: number
  emotional: number
}

const MS_PER_DAY = 24 * 60 * 60 * 1000

function toDate(value: string): Date {
  return new Date(`${value}T00:00:00`)
}

export function addDays(date: string, days: number): string {
  const next = toDate(date)
  next.setDate(next.getDate() + days)
  return next.toLocaleDateString('en-CA')
}

export function daysBetween(from: string, to: string): number {
  const start = toDate(from)
  const end = toDate(to)
  return Math.round((end.getTime() - start.getTime()) / MS_PER_DAY)
}

export function getDaysBetween(from: string, to: string): string[] {
  const totalDays = daysBetween(from, to)
  if (totalDays < 0) {
    return []
  }
  const dates: string[] = []
  for (let offset = 0; offset <= totalDays; offset += 1) {
    dates.push(addDays(from, offset))
  }
  return dates
}

export function formatDisplayDate(date: string): string {
  return new Date(`${date}T00:00:00`).toLocaleDateString('en-US', {
    month: 'long',
    day: 'numeric',
    year: 'numeric',
  })
}

export async function getEnergyLevels(
  date: string,
  token: string,
  signal?: AbortSignal,
): Promise<EnergyLevels | null> {
  const response = await fetch(
    `${API_BASE}/energy/levels?date=${encodeURIComponent(date)}`,
    {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      signal,
    },
  )

  if (response.status === 404) {
    return null
  }

  if (!response.ok) {
    throw new Error('failed_to_fetch_energy_levels')
  }

  const payload: EnergyLevels = await response.json()
  return payload
}

export async function getEnergyLevelsRange(
  from: string,
  to: string,
  token: string,
  signal?: AbortSignal,
): Promise<EnergyLevels[]> {
  const response = await fetch(
    `${API_BASE}/energy/levels/range?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}`,
    {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      signal,
    },
  )

  if (!response.ok) {
    throw new Error('failed_to_fetch_energy_levels_range')
  }

  const payload: EnergyLevels[] = await response.json()
  return payload
}

export async function saveEnergyLevels(
  levels: EnergyLevels,
  token: string,
): Promise<EnergyLevels> {
  const response = await fetch(`${API_BASE}/energy/levels`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(levels),
  })

  if (!response.ok) {
    throw new Error('failed_to_save_energy_levels')
  }

  const payload: EnergyLevels = await response.json()
  return payload
}
