import type { EnergyLevels } from '@/services/energyLevels'

export const ENERGY_LEVELS_RANGE_CACHE_KEY = 'energy-levels:range-cache'
export const ENERGY_LEVELS_FORCE_REFRESH_KEY = 'energy-levels:force-refresh'

export interface EnergyLevelsRangeCache {
  from: string
  to: string
  levels: EnergyLevels[]
  status: 'success' | 'empty'
}

export function readEnergyLevelsRangeCache(): EnergyLevelsRangeCache | null {
  if (typeof window === 'undefined') {
    return null
  }

  const raw = window.sessionStorage.getItem(ENERGY_LEVELS_RANGE_CACHE_KEY)
  if (!raw) {
    return null
  }

  try {
    const parsed = JSON.parse(raw) as EnergyLevelsRangeCache
    if (!parsed.from || !parsed.to || !Array.isArray(parsed.levels)) {
      return null
    }
    if (parsed.status !== 'success' && parsed.status !== 'empty') {
      return null
    }
    return parsed
  } catch {
    return null
  }
}

export function writeEnergyLevelsRangeCache(cache: EnergyLevelsRangeCache): void {
  if (typeof window === 'undefined') {
    return
  }

  window.sessionStorage.setItem(
    ENERGY_LEVELS_RANGE_CACHE_KEY,
    JSON.stringify(cache),
  )
}

export function clearEnergyLevelsRangeCache(): void {
  if (typeof window === 'undefined') {
    return
  }

  window.sessionStorage.removeItem(ENERGY_LEVELS_RANGE_CACHE_KEY)
}
