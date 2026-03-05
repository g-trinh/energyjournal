import { act, fireEvent, render, screen, waitFor } from '@testing-library/react'
import type { ReactNode } from 'react'
import { MemoryRouter } from 'react-router-dom'

import App from '@/App'
import { getCalendarStatus } from '@/services/calendar'

vi.mock('recharts', () => ({
  ResponsiveContainer: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  BarChart: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  Bar: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  XAxis: () => null,
  YAxis: () => null,
  Tooltip: () => null,
  Cell: () => null,
}))

vi.mock('@/services/calendar', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/services/calendar')>()
  return {
    ...actual,
    getCalendarStatus: vi.fn(),
    getCalendars: vi.fn(),
    getCalendarAuthURL: vi.fn(),
    saveCalendarConnection: vi.fn(),
  }
})

vi.mock('@/lib/session', () => ({
  getIdToken: () => 'token',
  clearSession: vi.fn(),
}))

function createSpendingsResponse(payload: Record<string, number>) {
  return {
    ok: true,
    status: 200,
    json: async () => payload,
  } as Response
}

describe('App time spending interactions', () => {
  beforeEach(() => {
    vi.useRealTimers()
    vi.mocked(getCalendarStatus).mockResolvedValue('connected')
  })

  it('requests previous week when clicking Prev week', async () => {
    const fetchMock = vi.fn().mockResolvedValue(createSpendingsResponse({ Tomato: 3 }))
    vi.stubGlobal('fetch', fetchMock)

    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>,
    )

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(1)
    })
    const firstRequestURL = new URL(String(fetchMock.mock.calls[0][0]), 'http://localhost')
    const firstStartDate = firstRequestURL.searchParams.get('start')

    fireEvent.click(screen.getByRole('button', { name: /previous week/i }))

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(2)
    })
    const secondRequestURL = new URL(String(fetchMock.mock.calls[1][0]), 'http://localhost')
    const secondStartDate = secondRequestURL.searchParams.get('start')

    expect(firstStartDate).not.toBeNull()
    expect(secondStartDate).not.toBeNull()

    const firstDate = new Date(String(firstStartDate))
    const secondDate = new Date(String(secondStartDate))
    const diffDays = (firstDate.getTime() - secondDate.getTime()) / (1000 * 60 * 60 * 24)
    expect(diffDays).toBe(7)
  })

  it('disables Next week on current week and enables it on past weeks', async () => {
    const fetchMock = vi.fn().mockResolvedValue(createSpendingsResponse({ Tomato: 3 }))
    vi.stubGlobal('fetch', fetchMock)

    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>,
    )

    const nextWeekButton = await screen.findByRole('button', { name: /next week/i })
    expect(nextWeekButton).toBeDisabled()

    fireEvent.click(screen.getByRole('button', { name: /previous week/i }))

    await waitFor(() => {
      expect(nextWeekButton).not.toBeDisabled()
    })
  })

  it('does not show spinner when loading finishes before 200ms', async () => {
    vi.useFakeTimers()
    const fetchMock = vi.fn().mockResolvedValue(createSpendingsResponse({ Tomato: 3 }))
    vi.stubGlobal('fetch', fetchMock)

    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>,
    )

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalled()
    })
    await act(async () => {
      vi.advanceTimersByTime(199)
    })

    expect(screen.queryByText('Gathering your energy data...')).not.toBeInTheDocument()
  })

  it('shows spinner after 200ms while loading and clears it when complete', async () => {
    vi.useFakeTimers()

    let resolveFetch: ((value: Response) => void) | null = null
    const fetchMock = vi.fn().mockImplementation(
      () =>
        new Promise<Response>((resolve) => {
          resolveFetch = resolve
        }),
    )
    vi.stubGlobal('fetch', fetchMock)

    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>,
    )

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalled()
    })

    await act(async () => {
      vi.advanceTimersByTime(201)
    })
    expect(screen.getByText('Gathering your energy data...')).toBeInTheDocument()

    await act(async () => {
      resolveFetch?.(createSpendingsResponse({ Tomato: 3 }))
    })
    await waitFor(() => {
      expect(screen.queryByText('Gathering your energy data...')).not.toBeInTheDocument()
    })
  })
})
