import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import type { ReactNode } from 'react'
import EnergyLevelsEditPage from './EnergyLevelsEditPage'
import {
  getEnergyLevels,
  saveEnergyLevels,
  type EnergyLevels,
} from '@/services/energyLevels'

vi.mock('react-dom', async (importOriginal) => {
  const original = await importOriginal<typeof import('react-dom')>()
  return {
    ...original,
    createPortal: (node: ReactNode) => node,
  }
})

vi.mock('@/services/energyLevels', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/services/energyLevels')>()
  return {
    ...actual,
    getEnergyLevels: vi.fn(),
    saveEnergyLevels: vi.fn(),
  }
})

vi.mock('@/lib/session', () => ({
  getIdToken: () => 'test-token',
}))

describe('EnergyLevelsEditPage', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('loads 404 state then saves and shows confirmation toast', async () => {
    vi.mocked(getEnergyLevels).mockResolvedValue(null)

    const saved: EnergyLevels = {
      date: '2026-02-21',
      physical: 5,
      mental: 5,
      emotional: 5,
    }
    vi.mocked(saveEnergyLevels).mockResolvedValue(saved)

    render(<EnergyLevelsEditPage />)

    const physicalSlider = await screen.findByRole('slider', {
      name: 'Physical energy level',
    })
    const mentalSlider = screen.getByRole('slider', { name: 'Mental energy level' })
    const emotionalSlider = screen.getByRole('slider', {
      name: 'Emotional energy level',
    })

    expect(physicalSlider).toHaveAttribute('aria-valuenow', '5')
    expect(mentalSlider).toHaveAttribute('aria-valuenow', '5')
    expect(emotionalSlider).toHaveAttribute('aria-valuenow', '5')

    fireEvent.click(screen.getByRole('button', { name: 'Save Energy Levels' }))

    await waitFor(() => {
      expect(saveEnergyLevels).toHaveBeenCalledWith(
        {
          date: expect.stringMatching(/^\d{4}-\d{2}-\d{2}$/),
          physical: 5,
          mental: 5,
          emotional: 5,
        },
        'test-token',
      )
    })

    expect(await screen.findByText('Energy levels saved')).toBeInTheDocument()
    expect(
      screen.getByText(/Physical 5 · Mental 5 · Emotional 5/),
    ).toBeInTheDocument()
  })
})
