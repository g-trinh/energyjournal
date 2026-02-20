import { fireEvent, render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import UserMenu from './UserMenu'
import { useAuth } from '@/contexts/AuthContext'

vi.mock('@/contexts/AuthContext', () => ({
  useAuth: vi.fn(),
}))

describe('UserMenu', () => {
  const mockedUseAuth = vi.mocked(useAuth)

  beforeEach(() => {
    mockedUseAuth.mockReset()
  })

  it('renders a skeleton while auth state is loading', () => {
    mockedUseAuth.mockReturnValue({
      status: 'loading',
      email: null,
      isAuthenticated: false,
      signOut: vi.fn(),
    })

    render(
      <MemoryRouter>
        <UserMenu />
      </MemoryRouter>,
    )

    expect(screen.getByLabelText('Loading user menu')).toBeInTheDocument()
  })

  it('navigates to /auth in anonymous mode', () => {
    mockedUseAuth.mockReturnValue({
      status: 'anonymous',
      email: null,
      isAuthenticated: false,
      signOut: vi.fn(),
    })

    render(
      <MemoryRouter initialEntries={['/']}>
        <Routes>
          <Route path="/" element={<UserMenu />} />
          <Route path="/auth" element={<h1>Auth Destination</h1>} />
        </Routes>
      </MemoryRouter>,
    )

    fireEvent.click(screen.getByRole('button', { name: 'Log in' }))
    expect(screen.getByRole('heading', { name: 'Auth Destination' })).toBeInTheDocument()
  })

  it('opens and closes dropdown, then signs out', () => {
    const signOut = vi.fn()
    mockedUseAuth.mockReturnValue({
      status: 'authenticated',
      email: 'person@example.com',
      isAuthenticated: true,
      signOut,
    })

    render(
      <MemoryRouter>
        <UserMenu />
      </MemoryRouter>,
    )

    fireEvent.click(screen.getByRole('button', { name: 'person@example.com' }))
    expect(screen.getByRole('menu')).toBeInTheDocument()

    fireEvent.keyDown(document, { key: 'Escape' })
    expect(screen.queryByRole('menu')).not.toBeInTheDocument()

    fireEvent.click(screen.getByRole('button', { name: 'person@example.com' }))
    fireEvent.click(screen.getByRole('menuitem', { name: 'Sign out' }))
    expect(signOut).toHaveBeenCalledTimes(1)
  })
})
