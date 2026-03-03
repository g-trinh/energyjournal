import { render, screen } from '@testing-library/react'
import { act } from 'react'
import Toast from './Toast'

describe('Toast', () => {
  it('auto dismisses after duration', () => {
    vi.useFakeTimers()
    const onDismiss = vi.fn()

    render(
      <Toast
        message="Calendar connected!"
        subtitle="Your events are syncing…"
        duration={4000}
        onDismiss={onDismiss}
      />,
    )

    expect(screen.getByRole('status')).toBeInTheDocument()

    act(() => {
      vi.advanceTimersByTime(4000)
    })

    expect(onDismiss).toHaveBeenCalledTimes(1)
    vi.useRealTimers()
  })
})
