import { fireEvent, render, screen } from '@testing-library/react'
import StepIndicator from './StepIndicator'

describe('StepIndicator', () => {
  it('renders correct aria states for step 1 active', () => {
    render(<StepIndicator currentStep={1} onStepClick={vi.fn()} />)

    expect(screen.getByRole('button', { name: 'Go to step 1' })).toHaveAttribute(
      'aria-current',
      'step',
    )
    expect(screen.getByRole('button', { name: 'Go to step 2' })).toBeDisabled()
  })

  it('renders correct aria states for step 2 active', () => {
    render(<StepIndicator currentStep={2} onStepClick={vi.fn()} />)

    expect(screen.getByRole('button', { name: 'Go to step 1' })).not.toHaveAttribute(
      'aria-current',
    )
    expect(screen.getByRole('button', { name: 'Go to step 2' })).toHaveAttribute(
      'aria-current',
      'step',
    )
  })

  it('fires onStepClick with the correct step number', () => {
    const onStepClick = vi.fn()
    render(<StepIndicator currentStep={2} onStepClick={onStepClick} />)

    fireEvent.click(screen.getByRole('button', { name: 'Go to step 1' }))
    fireEvent.click(screen.getByRole('button', { name: 'Go to step 2' }))

    expect(onStepClick).toHaveBeenCalledWith(1)
    expect(onStepClick).toHaveBeenCalledWith(2)
  })
})
