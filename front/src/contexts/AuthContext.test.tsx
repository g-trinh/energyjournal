import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { AuthProvider, useAuth } from './AuthContext'
import { clearSession, getIdToken, isAuthenticated } from '@/lib/session'

const mockedNavigate = vi.fn()

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>(
    'react-router-dom',
  )

  return {
    ...actual,
    useNavigate: () => mockedNavigate,
  }
})

vi.mock('@/lib/session', () => ({
  clearSession: vi.fn(),
  getIdToken: vi.fn(),
  isAuthenticated: vi.fn(),
}))

function Probe() {
  const { status, email, signOut } = useAuth()
  return (
    <div>
      <p>Status: {status}</p>
      <p>Email: {email ?? 'none'}</p>
      <button type="button" onClick={signOut}>
        Sign out now
      </button>
    </div>
  )
}

describe('AuthContext', () => {
  const mockedIsAuthenticated = vi.mocked(isAuthenticated)
  const mockedGetIdToken = vi.mocked(getIdToken)
  const mockedClearSession = vi.mocked(clearSession)

  beforeEach(() => {
    mockedNavigate.mockReset()
    mockedClearSession.mockReset()
    mockedIsAuthenticated.mockReset()
    mockedGetIdToken.mockReset()
    vi.stubGlobal('fetch', vi.fn())
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('loads email and marks status authenticated on successful /users/me', async () => {
    mockedIsAuthenticated.mockReturnValue(true)
    mockedGetIdToken.mockReturnValue('id-token')
    vi.mocked(fetch).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({ email: 'hello@example.com' }),
    } as Response)

    render(
      <MemoryRouter>
        <AuthProvider>
          <Probe />
        </AuthProvider>
      </MemoryRouter>,
    )

    await waitFor(() =>
      expect(screen.getByText('Status: authenticated')).toBeInTheDocument(),
    )
    expect(screen.getByText('Email: hello@example.com')).toBeInTheDocument()
  })

  it('clears session and falls back to anonymous on 401', async () => {
    mockedIsAuthenticated.mockReturnValue(true)
    mockedGetIdToken.mockReturnValue('id-token')
    vi.mocked(fetch).mockResolvedValueOnce({
      ok: false,
      status: 401,
      json: async () => ({}),
    } as Response)

    render(
      <MemoryRouter>
        <AuthProvider>
          <Probe />
        </AuthProvider>
      </MemoryRouter>,
    )

    await waitFor(() =>
      expect(screen.getByText('Status: anonymous')).toBeInTheDocument(),
    )
    expect(mockedClearSession).toHaveBeenCalledTimes(1)
  })

  it('falls back to anonymous on network errors', async () => {
    mockedIsAuthenticated.mockReturnValue(true)
    mockedGetIdToken.mockReturnValue('id-token')
    vi.mocked(fetch).mockRejectedValueOnce(new Error('network'))

    render(
      <MemoryRouter>
        <AuthProvider>
          <Probe />
        </AuthProvider>
      </MemoryRouter>,
    )

    await waitFor(() =>
      expect(screen.getByText('Status: anonymous')).toBeInTheDocument(),
    )
    expect(mockedClearSession).not.toHaveBeenCalled()
  })

  it('signOut clears session and navigates to /auth', async () => {
    mockedIsAuthenticated.mockReturnValue(false)
    mockedGetIdToken.mockReturnValue(null)

    render(
      <MemoryRouter>
        <AuthProvider>
          <Probe />
        </AuthProvider>
      </MemoryRouter>,
    )

    await waitFor(() =>
      expect(screen.getByText('Status: anonymous')).toBeInTheDocument(),
    )

    fireEvent.click(screen.getByRole('button', { name: 'Sign out now' }))

    expect(mockedClearSession).toHaveBeenCalledTimes(1)
    expect(mockedNavigate).toHaveBeenCalledWith('/auth', { replace: true })
  })

  it('throws when useAuth is used outside AuthProvider', () => {
    expect(() => render(<Probe />)).toThrow(
      'useAuth must be used within an AuthProvider',
    )
  })
})
