const API_BASE = import.meta.env.VITE_API_BASE_URL ?? '/api'

export interface EnergyLevels {
  date: string
  physical: number
  mental: number
  emotional: number
  sleepQuality?: number
  stressLevel?: number
  physicalActivity?: string
  nutrition?: string
  socialInteractions?: string
  timeOutdoors?: string
  notes?: string
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

  const payload = await response.json()
  return parseEnergyLevels(payload)
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

  const payload: unknown = await response.json()
  if (!Array.isArray(payload)) {
    return []
  }
  return payload.map(parseEnergyLevels)
}

export async function saveEnergyLevels(
  levels: EnergyLevels,
  token: string,
): Promise<EnergyLevels> {
  const body: EnergyLevels = {
    date: levels.date,
    physical: levels.physical,
    mental: levels.mental,
    emotional: levels.emotional,
  }

  if (levels.sleepQuality !== undefined) {
    body.sleepQuality = levels.sleepQuality
  }
  if (levels.stressLevel !== undefined) {
    body.stressLevel = levels.stressLevel
  }
  if (levels.physicalActivity) {
    body.physicalActivity = levels.physicalActivity
  }
  if (levels.nutrition) {
    body.nutrition = levels.nutrition
  }
  if (levels.socialInteractions) {
    body.socialInteractions = levels.socialInteractions
  }
  if (levels.timeOutdoors) {
    body.timeOutdoors = levels.timeOutdoors
  }
  if (levels.notes) {
    body.notes = levels.notes
  }

  const response = await fetch(`${API_BASE}/energy/levels`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(body),
  })

  if (!response.ok) {
    throw new Error('failed_to_save_energy_levels')
  }

  const payload = await response.json()
  return parseEnergyLevels(payload)
}

function parseEnergyLevels(payload: unknown): EnergyLevels {
  const data = payload as Partial<EnergyLevels>

  return {
    date: data.date ?? '',
    physical: data.physical ?? 0,
    mental: data.mental ?? 0,
    emotional: data.emotional ?? 0,
    sleepQuality: normalizeOptionalNumber(data.sleepQuality),
    stressLevel: normalizeOptionalNumber(data.stressLevel),
    physicalActivity: normalizeOptionalString(data.physicalActivity),
    nutrition: normalizeOptionalString(data.nutrition),
    socialInteractions: normalizeOptionalString(data.socialInteractions),
    timeOutdoors: normalizeOptionalString(data.timeOutdoors),
    notes: normalizeOptionalString(data.notes),
  }
}

function normalizeOptionalNumber(value: unknown): number | undefined {
  if (typeof value !== 'number') {
    return undefined
  }
  return value
}

function normalizeOptionalString(value: unknown): string | undefined {
  if (typeof value !== 'string' || value === '') {
    return undefined
  }
  return value
}
