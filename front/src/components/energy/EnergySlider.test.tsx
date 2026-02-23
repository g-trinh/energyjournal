import { fireEvent, render, screen } from '@testing-library/react'
import EnergySlider from './EnergySlider'

describe('EnergySlider', () => {
  it('renders range input with expected attributes', () => {
    render(
      <EnergySlider
        value={5}
        color="#c4826d"
        ariaLabel="Physical energy level"
        onChange={vi.fn()}
      />,
    )

    const slider = screen.getByRole('slider', { name: 'Physical energy level' })
    expect(slider).toHaveAttribute('type', 'range')
    expect(slider).toHaveAttribute('min', '0')
    expect(slider).toHaveAttribute('max', '10')
    expect(slider).toHaveAttribute('step', '1')
    expect(slider).toHaveAttribute('aria-valuenow', '5')
  })

  it('calls onChange with parsed integer value', () => {
    const onChange = vi.fn()

    render(
      <EnergySlider
        value={5}
        color="#7eb8b3"
        ariaLabel="Mental energy level"
        onChange={onChange}
      />,
    )

    const slider = screen.getByRole('slider', { name: 'Mental energy level' })
    fireEvent.change(slider, { target: { value: '9' } })

    expect(onChange).toHaveBeenCalledWith(9)
  })
})
