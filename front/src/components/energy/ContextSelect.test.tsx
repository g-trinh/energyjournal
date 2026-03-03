import { fireEvent, render, screen } from '@testing-library/react'
import ContextSelect from './ContextSelect'

const OPTIONS = [
  { value: '', label: '—' },
  { value: 'none', label: 'None' },
  { value: 'light', label: 'Light' },
]

describe('ContextSelect', () => {
  it('opens on click and closes on outside click and escape', () => {
    render(
      <ContextSelect
        label="Physical Activity"
        options={OPTIONS}
        value=""
        onChange={vi.fn()}
      />,
    )

    const trigger = screen.getByRole('combobox', { name: 'Physical Activity' })
    fireEvent.click(trigger)
    expect(screen.getByRole('listbox', { name: 'Physical Activity' })).toBeInTheDocument()

    fireEvent.keyDown(trigger, { key: 'Escape' })
    expect(
      screen.queryByRole('listbox', { name: 'Physical Activity' }),
    ).not.toBeInTheDocument()

    fireEvent.click(trigger)
    expect(screen.getByRole('listbox', { name: 'Physical Activity' })).toBeInTheDocument()

    fireEvent.mouseDown(document.body)
    expect(
      screen.queryByRole('listbox', { name: 'Physical Activity' }),
    ).not.toBeInTheDocument()
  })

  it('supports keyboard navigation and selection', () => {
    const onChange = vi.fn()

    render(
      <ContextSelect
        label="Physical Activity"
        options={OPTIONS}
        value=""
        onChange={onChange}
      />,
    )

    const trigger = screen.getByRole('combobox', { name: 'Physical Activity' })
    fireEvent.keyDown(trigger, { key: 'ArrowDown' })

    const listbox = screen.getByRole('listbox', { name: 'Physical Activity' })
    fireEvent.keyDown(listbox, { key: 'ArrowDown' })
    fireEvent.keyDown(listbox, { key: 'Enter' })

    expect(onChange).toHaveBeenCalledWith('none')
  })

  it('calls onChange when selecting an option with click', () => {
    const onChange = vi.fn()

    render(
      <ContextSelect
        label="Nutrition"
        options={OPTIONS}
        value=""
        onChange={onChange}
      />,
    )

    fireEvent.click(screen.getByRole('combobox', { name: 'Nutrition' }))
    fireEvent.click(screen.getByRole('option', { name: 'Light' }))

    expect(onChange).toHaveBeenCalledWith('light')
  })

  it('does not open when disabled', () => {
    render(
      <ContextSelect
        label="Social Interactions"
        options={OPTIONS}
        value=""
        onChange={vi.fn()}
        disabled
      />,
    )

    const trigger = screen.getByRole('combobox', { name: 'Social Interactions' })
    expect(trigger).toBeDisabled()

    fireEvent.click(trigger)
    expect(
      screen.queryByRole('listbox', { name: 'Social Interactions' }),
    ).not.toBeInTheDocument()
  })
})
