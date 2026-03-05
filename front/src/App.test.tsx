import { render, screen, waitFor } from '@testing-library/react'
import type { ReactNode } from 'react'
import { MemoryRouter } from 'react-router-dom'
import App from './App'
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

describe('App calendar status gate', () => {
  it('renders connect prompt when status is disconnected', async () => {
    vi.mocked(getCalendarStatus).mockResolvedValue('disconnected')

    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>,
    )

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Connect your Google Calendar' })).toBeInTheDocument()
    })
    expect(screen.queryByRole('heading', { name: 'Time Distribution' })).not.toBeInTheDocument()
  })

  it('renders chart successfully on mobile viewport', async () => {
    vi.mocked(getCalendarStatus).mockResolvedValue('connected')
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({
        BlueberryCalendarLongName: 8,
      }),
    } as Response))

    Object.defineProperty(window, 'innerWidth', { configurable: true, value: 480 })
    window.dispatchEvent(new Event('resize'))

    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>,
    )

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Time Distribution' })).toBeInTheDocument()
    })
    expect(screen.getByRole('button', { name: /previous week/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /next week/i })).toBeInTheDocument()
  })
})
