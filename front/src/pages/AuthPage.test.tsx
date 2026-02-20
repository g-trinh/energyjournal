import { fireEvent, render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import AuthPage from './AuthPage'
import { useAuth } from '@/contexts/AuthContext'

vi.mock('@/components/auth/LoginCard', () => ({
  default: ({ onLoginSuccess }: { onLoginSuccess: (tokens: { idToken: string; refreshToken: string; expiresIn: string }) => void }) => (
    <button
      type="button"
      onClick={() =>
        onLoginSuccess({
          idToken: 'id-token',
          refreshToken: 'refresh-token',
          expiresIn: '3600',
        })
      }
    >
      Mock Login
    </button>
  ),
}))

vi.mock('@/components/auth/SignupCard', () => ({
  default: () => <div>Signup Card</div>,
}))

vi.mock('@/contexts/AuthContext', () => ({
  useAuth: vi.fn(),
}))

describe('AuthPage navigation', () => {
  const mockedUseAuth = vi.mocked(useAuth)

  beforeEach(() => {
    localStorage.clear()
    mockedUseAuth.mockReset()
    mockedUseAuth.mockReturnValue({
      status: 'anonymous',
      email: null,
      isAuthenticated: false,
      signIn: vi.fn(),
      signOut: vi.fn(),
    })
  })

  it('stores tokens and redirects to /timespending after login success', () => {
    render(
      <MemoryRouter initialEntries={['/auth']}>
        <Routes>
          <Route path="/auth" element={<AuthPage />} />
          <Route path="/timespending" element={<h1>Timespending Destination</h1>} />
        </Routes>
      </MemoryRouter>,
    )

    fireEvent.click(screen.getByRole('button', { name: 'Mock Login' }))

    expect(
      screen.getByRole('heading', { name: 'Timespending Destination' }),
    ).toBeInTheDocument()
    expect(localStorage.getItem('idToken')).toBe('id-token')
    expect(localStorage.getItem('refreshToken')).toBe('refresh-token')
  })
})
