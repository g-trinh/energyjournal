import { formatDisplayDate } from '@/services/energyLevels'
import { ENERGY_COLORS, ENERGY_LABELS, type EnergyDimension } from '@/lib/energyColors'

interface EnergyTooltipPayload {
  dataKey?: EnergyDimension
  value?: number | null
  color?: string
}

interface EnergyTooltipProps {
  active?: boolean
  payload?: EnergyTooltipPayload[]
  label?: string
}

export default function EnergyTooltip({ active, payload, label }: EnergyTooltipProps) {
  if (!active || !payload || payload.length === 0 || !label) {
    return null
  }

  const items = (Object.keys(ENERGY_COLORS) as EnergyDimension[]).map((dimension) => {
    const entry = payload.find((item) => item.dataKey === dimension)
    return {
      dimension,
      value: entry?.value ?? null,
      color: entry?.color ?? ENERGY_COLORS[dimension],
      label: ENERGY_LABELS[dimension],
    }
  })

  return (
    <div className="energy-levels-tooltip">
      <p className="energy-levels-tooltip-title">{formatDisplayDate(label)}</p>
      <div className="energy-levels-tooltip-divider" aria-hidden="true" />
      <div className="energy-levels-tooltip-rows">
        {items.map((item) => (
          <div key={item.dimension} className="energy-levels-tooltip-row">
            <span className="energy-levels-tooltip-label">
              <span
                className="energy-levels-tooltip-dot"
                style={{ backgroundColor: item.color }}
                aria-hidden="true"
              />
              {item.label}
            </span>
            <span className="energy-levels-tooltip-value">
              {item.value ?? '--'}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}
