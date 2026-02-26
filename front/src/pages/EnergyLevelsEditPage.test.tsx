import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import type { ReactNode } from 'react'
import { MemoryRouter } from 'react-router-dom'
import EnergyLevelsEditPage from './EnergyLevelsEditPage'
import {
  ENERGY_LEVELS_FORCE_REFRESH_KEY,
  ENERGY_LEVELS_RANGE_CACHE_KEY,
} from '@/lib/energyLevelsCache'
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

const navigateMock = vi.fn()
vi.mock('react-router-dom', async (importOriginal) => {
  const original = await importOriginal<typeof import('react-router-dom')>()
  return {
    ...original,
    useNavigate: () => navigateMock,
  }
})

describe('EnergyLevelsEditPage', () => {
  beforeEach(() => {
    sessionStorage.clear()
    navigateMock.mockReset()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('supports full step flow with existing context values and saves full payload', async () => {
    vi.mocked(getEnergyLevels).mockResolvedValue({
      date: '2026-02-21',
      physical: 7,
      mental: 6,
      emotional: 8,
      sleepQuality: 4,
      stressLevel: 2,
      physicalActivity: 'moderate',
      nutrition: 'good',
      socialInteractions: 'positive',
      timeOutdoors: '30min_1hr',
      notes: 'Felt focused.',
    })

    const saved: EnergyLevels = {
      date: '2026-02-21',
      physical: 7,
      mental: 6,
      emotional: 8,
      sleepQuality: 5,
      stressLevel: 1,
      physicalActivity: 'intense',
      nutrition: 'excellent',
      socialInteractions: 'neutral',
      timeOutdoors: 'over_1hr',
      notes: 'Finished strong.',
    }
    vi.mocked(saveEnergyLevels).mockResolvedValue(saved)

    sessionStorage.setItem(
      ENERGY_LEVELS_RANGE_CACHE_KEY,
      '{"from":"2026-02-10","to":"2026-02-23","levels":[],"status":"empty"}',
    )

    render(
      <MemoryRouter>
        <EnergyLevelsEditPage />
      </MemoryRouter>,
    )

    const physicalSlider = await screen.findByRole('slider', {
      name: 'Physical energy level',
    })
    expect(physicalSlider).toHaveAttribute('aria-valuenow', '7')

    fireEvent.click(screen.getByRole('button', { name: /next step/i }))

    expect(await screen.findByText('Daily Context')).toBeInTheDocument()
    expect(
      screen.getByRole('slider', { name: 'Sleep Quality energy level' }),
    ).toHaveAttribute('aria-valuenow', '4')
    expect(
      screen.getByRole('slider', { name: 'Stress Level energy level' }),
    ).toHaveAttribute('aria-valuenow', '2')
    expect(
      screen.getByRole('combobox', { name: 'Physical Activity' }),
    ).toHaveTextContent('Moderate')
    expect(screen.getByRole('combobox', { name: 'Nutrition' })).toHaveTextContent(
      'Good quality',
    )
    expect(
      screen.getByRole('combobox', { name: 'Social Interactions' }),
    ).toHaveTextContent('Positive')
    expect(
      screen.getByRole('combobox', { name: 'Time Outdoors' }),
    ).toHaveTextContent('30 min-1 hr')
    expect(screen.getByLabelText('Daily Reflection')).toHaveValue('Felt focused.')

    fireEvent.click(screen.getByRole('button', { name: '← Previous Step' }))
    expect(
      screen.getByRole('slider', { name: 'Physical energy level' }),
    ).toHaveAttribute('aria-valuenow', '7')

    fireEvent.click(screen.getByRole('button', { name: /next step/i }))

    fireEvent.change(
      screen.getByRole('slider', { name: 'Sleep Quality energy level' }),
      { target: { value: '5' } },
    )
    fireEvent.change(
      screen.getByRole('slider', { name: 'Stress Level energy level' }),
      { target: { value: '1' } },
    )

    fireEvent.click(screen.getByRole('combobox', { name: 'Physical Activity' }))
    fireEvent.click(screen.getByRole('option', { name: 'Intense' }))

    fireEvent.click(screen.getByRole('combobox', { name: 'Nutrition' }))
    fireEvent.click(screen.getByRole('option', { name: 'Excellent quality' }))

    fireEvent.click(screen.getByRole('combobox', { name: 'Social Interactions' }))
    fireEvent.click(screen.getByRole('option', { name: 'Neutral' }))

    fireEvent.click(screen.getByRole('combobox', { name: 'Time Outdoors' }))
    fireEvent.click(screen.getByRole('option', { name: 'Over 1 hr' }))

    fireEvent.change(screen.getByLabelText('Daily Reflection'), {
      target: { value: 'Finished strong.' },
    })

    fireEvent.click(screen.getByRole('button', { name: 'Save Entry' }))

    await waitFor(() => {
      expect(saveEnergyLevels).toHaveBeenCalledWith(
        {
          date: expect.stringMatching(/^\d{4}-\d{2}-\d{2}$/),
          physical: 7,
          mental: 6,
          emotional: 8,
          sleepQuality: 5,
          stressLevel: 1,
          physicalActivity: 'intense',
          nutrition: 'excellent',
          socialInteractions: 'neutral',
          timeOutdoors: 'over_1hr',
          notes: 'Finished strong.',
        },
        'test-token',
      )
    })

    expect(await screen.findByText('Energy levels saved')).toBeInTheDocument()
    expect(
      screen.getByText(/Physical 7 · Mental 6 · Emotional 8/),
    ).toBeInTheDocument()
    expect(sessionStorage.getItem(ENERGY_LEVELS_RANGE_CACHE_KEY)).toBeNull()
    expect(sessionStorage.getItem(ENERGY_LEVELS_FORCE_REFRESH_KEY)).toBe('1')
  })

  it('navigates to levels when clicking Back and Cancel on step 1', async () => {
    vi.mocked(getEnergyLevels).mockResolvedValue(null)

    render(
      <MemoryRouter>
        <EnergyLevelsEditPage />
      </MemoryRouter>,
    )

    fireEvent.click(await screen.findByRole('button', { name: /back to energy levels/i }))
    expect(navigateMock).toHaveBeenCalledWith('/energy/levels')

    fireEvent.click(screen.getByRole('button', { name: 'Cancel' }))
    expect(navigateMock).toHaveBeenCalledWith('/energy/levels')
  })
})
