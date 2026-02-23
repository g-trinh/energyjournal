import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import type { ReactNode } from 'react'
import { MemoryRouter } from 'react-router-dom'
import { ENERGY_LEVELS_FORCE_REFRESH_KEY, ENERGY_LEVELS_RANGE_CACHE_KEY } from '@/lib/energyLevelsCache'
import EnergyLevelsPage, { buildChartData } from './EnergyLevelsPage'
import { getEnergyLevelsRange } from '@/services/energyLevels'

vi.mock('recharts', () => ({
  ResponsiveContainer: ({ children }: { children: ReactNode }) => (
    <div data-testid="responsive-container">{children}</div>
  ),
  LineChart: ({ children, data }: { children: ReactNode; data?: unknown[] }) => (
    <div data-testid="line-chart" data-points={data?.length ?? 0}>
      {children}
    </div>
  ),
  Line: ({ dataKey, hide }: { dataKey: string; hide?: boolean }) => (
    <div data-testid={`line-${dataKey}`} data-hide={hide ? 'true' : 'false'} />
  ),
  XAxis: () => null,
  YAxis: () => null,
  CartesianGrid: () => null,
  Tooltip: ({ content }: { content: ReactNode }) => (
    <div data-testid="tooltip">{content}</div>
  ),
}))

vi.mock('@/services/energyLevels', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/services/energyLevels')>()
  return {
    ...actual,
    getEnergyLevelsRange: vi.fn(),
  }
})

vi.mock('@/lib/session', () => ({
  getIdToken: () => 'test-token',
}))

describe('buildChartData', () => {
  it('fills missing days with null values', () => {
    const data = buildChartData(
      [{ date: '2026-02-10', physical: 6, mental: 5, emotional: 4 }],
      '2026-02-09',
      '2026-02-11',
    )

    expect(data).toEqual([
      { date: '2026-02-09', physical: null, mental: null, emotional: null },
      { date: '2026-02-10', physical: 6, mental: 5, emotional: 4 },
      { date: '2026-02-11', physical: null, mental: null, emotional: null },
    ])
  })
})

describe('EnergyLevelsPage', () => {
  const originalMatchMedia = window.matchMedia

  beforeEach(() => {
    sessionStorage.clear()
    vi.useFakeTimers({ shouldAdvanceTime: true })
    vi.setSystemTime(new Date('2026-02-23T12:00:00'))
    Object.defineProperty(window, 'matchMedia', {
      writable: true,
      value: vi.fn().mockImplementation((query: string) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(),
        removeListener: vi.fn(),
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
      })),
    })
  })

  afterEach(() => {
    vi.useRealTimers()
    vi.restoreAllMocks()
    Object.defineProperty(window, 'matchMedia', {
      writable: true,
      value: originalMatchMedia,
    })
  })

  function renderPage() {
    return render(
      <MemoryRouter>
        <EnergyLevelsPage />
      </MemoryRouter>,
    )
  }

  it('shows loading state while fetching', () => {
    vi.mocked(getEnergyLevelsRange).mockImplementation(
      () => new Promise(() => {}) as Promise<never>,
    )

    renderPage()
    expect(screen.getByText('Gathering your energy data...')).toBeInTheDocument()
  })

  it('renders empty state on empty response', async () => {
    vi.mocked(getEnergyLevelsRange).mockResolvedValue([])

    renderPage()

    expect(await screen.findByText('No entries yet')).toBeInTheDocument()
  })

  it('renders error state on fetch failure', async () => {
    vi.mocked(getEnergyLevelsRange).mockRejectedValue(new Error('boom'))

    renderPage()

    expect(await screen.findByText('Unable to load data')).toBeInTheDocument()
  })

  it('renders chart and legend on success', async () => {
    vi.mocked(getEnergyLevelsRange).mockResolvedValue([
      { date: '2026-02-22', physical: 6, mental: 5, emotional: 4 },
      { date: '2026-02-23', physical: 7, mental: 6, emotional: 8 },
    ])

    renderPage()

    expect(await screen.findByTestId('line-physical')).toBeInTheDocument()
    expect(screen.getByTestId('line-mental')).toBeInTheDocument()
    expect(screen.getByTestId('line-emotional')).toBeInTheDocument()
  })

  it('updates date range when selecting presets', async () => {
    vi.mocked(getEnergyLevelsRange).mockResolvedValue([])

    renderPage()

    await waitFor(() => {
      expect(getEnergyLevelsRange).toHaveBeenCalled()
    })

    fireEvent.click(screen.getByRole('button', { name: '7d' }))

    await waitFor(() => {
      expect(getEnergyLevelsRange).toHaveBeenLastCalledWith(
        '2026-02-17',
        '2026-02-23',
        'test-token',
        expect.any(AbortSignal),
      )
    })
  })

  it('shows clamp warning when range exceeds 30 days', async () => {
    vi.mocked(getEnergyLevelsRange).mockResolvedValue([])

    renderPage()

    const fromInput = await screen.findByLabelText('FROM')
    fireEvent.change(fromInput, { target: { value: '2026-01-01' } })

    expect(await screen.findByText('Range limited to 30 days')).toBeInTheDocument()
  })

  it('toggles legend visibility state', async () => {
    vi.mocked(getEnergyLevelsRange).mockResolvedValue([
      { date: '2026-02-22', physical: 6, mental: 5, emotional: 4 },
    ])

    renderPage()

    const toggle = await screen.findByRole('button', { name: 'Toggle Physical line' })
    expect(screen.getByTestId('line-physical')).toHaveAttribute('data-hide', 'false')

    fireEvent.click(toggle)

    expect(toggle).toHaveAttribute('aria-pressed', 'true')
    expect(screen.getByTestId('line-physical')).toHaveAttribute('data-hide', 'true')
  })

  it('ignores cached range when refresh marker is present', async () => {
    sessionStorage.setItem(
      ENERGY_LEVELS_RANGE_CACHE_KEY,
      JSON.stringify({
        from: '2026-02-01',
        to: '2026-02-07',
        levels: [{ date: '2026-02-02', physical: 6, mental: 4, emotional: 5 }],
        status: 'success',
      }),
    )
    sessionStorage.setItem(ENERGY_LEVELS_FORCE_REFRESH_KEY, '1')
    vi.mocked(getEnergyLevelsRange).mockResolvedValue([])

    renderPage()

    await waitFor(() => {
      expect(getEnergyLevelsRange).toHaveBeenLastCalledWith(
        '2026-02-10',
        '2026-02-23',
        'test-token',
        expect.any(AbortSignal),
      )
    })
    expect(sessionStorage.getItem(ENERGY_LEVELS_FORCE_REFRESH_KEY)).toBeNull()
    expect(sessionStorage.getItem(ENERGY_LEVELS_RANGE_CACHE_KEY)).toContain(
      '"status":"empty"',
    )
  })
})
