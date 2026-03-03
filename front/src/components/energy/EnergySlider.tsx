import type { ChangeEvent, CSSProperties } from 'react'

interface EnergySliderProps {
  value: number
  color: string
  ariaLabel: string
  onChange: (value: number) => void
  disabled?: boolean
  min?: number
  max?: number
}

export default function EnergySlider({
  value,
  color,
  ariaLabel,
  onChange,
  disabled = false,
  min = 0,
  max = 10,
}: EnergySliderProps) {
  const range = max - min
  const style = {
    '--fill': `${range === 0 ? 0 : ((value - min) / range) * 100}%`,
    '--track-color': color,
  } as CSSProperties

  function handleChange(event: ChangeEvent<HTMLInputElement>) {
    onChange(parseInt(event.target.value, 10))
  }

  return (
    <input
      type="range"
      className="energy-slider"
      min={min}
      max={max}
      step={1}
      value={value}
      onChange={handleChange}
      disabled={disabled}
      aria-label={ariaLabel}
      aria-valuemin={min}
      aria-valuemax={max}
      aria-valuenow={value}
      style={style}
    />
  )
}
