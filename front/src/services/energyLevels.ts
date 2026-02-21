const API_BASE = import.meta.env.VITE_API_BASE_URL ?? '/api'

export interface EnergyLevels {
  date: string
  physical: number
  mental: number
  emotional: number
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
