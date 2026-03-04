import { fireEvent, render, screen } from '@testing-library/react'
import { act } from 'react'
import CalendarPicker from './CalendarPicker'

describe('CalendarPicker', () => {
  it('calls onSave with selected calendar id', async () => {
    const onSave = vi.fn().mockResolvedValue(undefined)

    render(
      <CalendarPicker
        calendars={[
          { id: 'primary', name: 'Primary', color: '#8fa58b' },
          { id: 'team', name: 'Team', color: '#4a6fa5' },
        ]}
        isLoading={false}
        onSave={onSave}
      />,
    )

    fireEvent.change(screen.getByLabelText('Select a calendar'), { target: { value: 'team' } })
    await act(async () => {
      fireEvent.click(screen.getByRole('button', { name: 'Save' }))
    })

    expect(onSave).toHaveBeenCalledWith('team')
  })
})
