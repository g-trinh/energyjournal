export const ENERGY_COLORS = {
  physical: '#c4826d',
  mental: '#7eb8b3',
  emotional: '#8fa58b',
} as const

export const ENERGY_LABELS: Record<EnergyDimension, string> = {
  physical: 'Physical',
  mental: 'Mental',
  emotional: 'Emotional',
}

export type EnergyDimension = keyof typeof ENERGY_COLORS
