import { fireEvent, render, screen } from '@testing-library/react'
import { act } from 'react'
import CalendarConnectPrompt from './CalendarConnectPrompt'

describe('CalendarConnectPrompt', () => {
  it('calls onConnect when the CTA is clicked', async () => {
    const onConnect = vi.fn().mockResolvedValue(undefined)
    render(<CalendarConnectPrompt onConnect={onConnect} />)

    await act(async () => {
      fireEvent.click(screen.getByRole('button', { name: 'Connect Google Calendar' }))
    })

    expect(onConnect).toHaveBeenCalledTimes(1)
  })
})
