import { fireEvent, render, screen } from '@testing-library/react'
import EnergySection from './EnergySection'

describe('EnergySection', () => {
  it('renders label and current value', () => {
    render(
      <EnergySection label="Physical" color="#c4826d" value={5} onChange={vi.fn()} />,
    )

    expect(screen.getByText('Physical')).toBeInTheDocument()
    expect(screen.getByText('5')).toBeInTheDocument()
    expect(screen.getByRole('slider', { name: 'Physical energy level' })).toHaveAttribute(
      'aria-valuenow',
      '5',
    )
  })

  it('shows divider by default and hides it when showDivider is false', () => {
    const { rerender } = render(
      <EnergySection label="Physical" color="#c4826d" value={5} onChange={vi.fn()} showDivider />,
    )
    expect(document.querySelector('.energy-section-divider')).toBeInTheDocument()

    rerender(
      <EnergySection label="Physical" color="#c4826d" value={5} onChange={vi.fn()} showDivider={false} />,
    )
    expect(document.querySelector('.energy-section-divider')).not.toBeInTheDocument()
  })

  it('calls onChange when the slider changes', () => {
    const onChange = vi.fn()
    render(
      <EnergySection label="Mental" color="#7eb8b3" value={3} onChange={onChange} />,
    )

    fireEvent.change(screen.getByRole('slider', { name: 'Mental energy level' }), {
      target: { value: '8' },
    })

    expect(onChange).toHaveBeenCalledWith(8)
  })
})
