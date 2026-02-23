import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import AppRouter from './AppRouter'

vi.mock('./App', () => ({
  default: () => <h1>Timespending Page</h1>,
}))

vi.mock('./pages/AuthPage', () => ({
  default: () => <h1>Auth Page</h1>,
}))

vi.mock('./pages/ActivatePage', () => ({
  default: () => <h1>Activate Page</h1>,
}))

vi.mock('./pages/LandingPage', () => ({
  default: () => <h1>Landing Page</h1>,
}))

vi.mock('./pages/EnergyLevelsEditPage', () => ({
  default: () => <h1>Energy Levels Page</h1>,
}))

vi.mock('./pages/EnergyLevelsPage', () => ({
  default: () => <h1>Energy Levels Chart Page</h1>,
}))

describe('AppRouter route access', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({ email: 'person@example.com' }),
      }),
    )
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('allows anonymous users to access landing, auth, and activate routes', async () => {
    const cases = [
      { path: '/', expected: 'Landing Page' },
      { path: '/auth', expected: 'Auth Page' },
      { path: '/activate', expected: 'Activate Page' },
    ]

    for (const testCase of cases) {
      const { unmount } = render(
        <MemoryRouter initialEntries={[testCase.path]}>
          <AuthProvider>
            <AppRouter />
          </AuthProvider>
        </MemoryRouter>,
      )
      expect(
        await screen.findByRole('heading', { name: testCase.expected }),
      ).toBeInTheDocument()
      expect(screen.getByText('Energy Journal')).toBeInTheDocument()
      unmount()
    }
  })

  it('redirects anonymous users from protected routes to /auth', async () => {
    const protectedCases = ['/timespending', '/energy/levels', '/energy/levels/edit']

    for (const path of protectedCases) {
      const { unmount } = render(
        <MemoryRouter initialEntries={[path]}>
          <AuthProvider>
            <AppRouter />
          </AuthProvider>
        </MemoryRouter>,
      )

      expect(
        await screen.findByRole('heading', { name: 'Auth Page' }),
      ).toBeInTheDocument()

      unmount()
    }
  })

  it('allows authenticated users to access protected routes', async () => {
    localStorage.setItem('idToken', 'token')

    const protectedCases = [
      { path: '/timespending', expected: 'Timespending Page' },
      { path: '/energy/levels', expected: 'Energy Levels Chart Page' },
      { path: '/energy/levels/edit', expected: 'Energy Levels Page' },
    ]

    for (const testCase of protectedCases) {
      const { unmount } = render(
        <MemoryRouter initialEntries={[testCase.path]}>
          <AuthProvider>
            <AppRouter />
          </AuthProvider>
        </MemoryRouter>,
      )

      expect(
        await screen.findByRole('heading', { name: testCase.expected }),
      ).toBeInTheDocument()

      unmount()
    }
  })

  it('redirects authenticated users from /auth and /activate to /timespending', async () => {
    localStorage.setItem('idToken', 'token')

    const protectedCases = ['/auth', '/activate']

    for (const path of protectedCases) {
      const { unmount } = render(
        <MemoryRouter initialEntries={[path]}>
          <AuthProvider>
            <AppRouter />
          </AuthProvider>
        </MemoryRouter>,
      )

      expect(
        await screen.findByRole('heading', { name: 'Timespending Page' }),
      ).toBeInTheDocument()

      unmount()
    }
  })
})
