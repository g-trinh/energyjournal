import EnergySlider from './EnergySlider'

interface EnergySectionProps {
  label: string
  color: string
  value: number
  onChange: (value: number) => void
  showDivider?: boolean
  disabled?: boolean
  min?: number
  max?: number
}

export default function EnergySection({
  label,
  color,
  value,
  onChange,
  showDivider = true,
  disabled = false,
  min = 0,
  max = 10,
}: EnergySectionProps) {
  return (
    <section className="energy-section">
      <div className="energy-section-header">
        <div className="energy-section-label-group">
          <span
            className="energy-section-dot"
            aria-hidden="true"
            style={{ backgroundColor: color }}
          />
          <span className="energy-section-label">{label}</span>
        </div>
        <div className="energy-section-value" aria-live="off">
          <span className="energy-section-value-number" style={{ color }}>
            {value}
          </span>
          <span className="energy-section-value-max">/ {max}</span>
        </div>
      </div>

      <EnergySlider
        value={value}
        color={color}
        ariaLabel={`${label} energy level`}
        onChange={onChange}
        disabled={disabled}
        min={min}
        max={max}
      />

      {showDivider && <hr className="energy-section-divider" />}
    </section>
  )
}
