import type { ChangeEvent, CSSProperties } from 'react'

interface EnergySliderProps {
  value: number
  color: string
  ariaLabel: string
  onChange: (value: number) => void
  disabled?: boolean
}

export default function EnergySlider({
  value,
  color,
  ariaLabel,
  onChange,
  disabled = false,
}: EnergySliderProps) {
  const style = {
    '--fill': `${(value / 10) * 100}%`,
    '--track-color': color,
  } as CSSProperties

  function handleChange(event: ChangeEvent<HTMLInputElement>) {
    onChange(parseInt(event.target.value, 10))
  }

  return (
    <input
      type="range"
      className="energy-slider"
      min={0}
      max={10}
      step={1}
      value={value}
      onChange={handleChange}
      disabled={disabled}
      aria-label={ariaLabel}
      aria-valuemin={0}
      aria-valuemax={10}
      aria-valuenow={value}
      style={style}
    />
  )
}
