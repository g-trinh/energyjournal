import { fireEvent, render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import BurgerMenu from './BurgerMenu'
import { useAuth } from '@/contexts/AuthContext'

vi.mock('@/contexts/AuthContext', () => ({
  useAuth: vi.fn(),
}))

describe('BurgerMenu', () => {
  const mockedUseAuth = vi.mocked(useAuth)

  beforeEach(() => {
    mockedUseAuth.mockReset()
  })

  it('renders account section only for authenticated users', () => {
    mockedUseAuth.mockReturnValue({
      status: 'authenticated',
      email: 'person@example.com',
      isAuthenticated: true,
      signOut: vi.fn(),
    })

    render(
      <MemoryRouter>
        <BurgerMenu
          open
          onClose={vi.fn()}
          triggerRef={{ current: document.createElement('button') }}
        />
      </MemoryRouter>,
    )

    expect(screen.getByText('ACCOUNT')).toBeInTheDocument()
    expect(screen.getByText('person@example.com')).toBeInTheDocument()
  })

  it('closes on Escape key', () => {
    const onClose = vi.fn()
    mockedUseAuth.mockReturnValue({
      status: 'anonymous',
      email: null,
      isAuthenticated: false,
      signOut: vi.fn(),
    })

    render(
      <MemoryRouter>
        <BurgerMenu
          open
          onClose={onClose}
          triggerRef={{ current: document.createElement('button') }}
        />
      </MemoryRouter>,
    )

    fireEvent.keyDown(document, { key: 'Escape' })
    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it('calls signOut from authenticated account row', () => {
    const onClose = vi.fn()
    const signOut = vi.fn()
    mockedUseAuth.mockReturnValue({
      status: 'authenticated',
      email: 'person@example.com',
      isAuthenticated: true,
      signOut,
    })

    render(
      <MemoryRouter>
        <BurgerMenu
          open
          onClose={onClose}
          triggerRef={{ current: document.createElement('button') }}
        />
      </MemoryRouter>,
    )

    fireEvent.click(screen.getByRole('button', { name: 'Sign out' }))
    expect(signOut).toHaveBeenCalledTimes(1)
    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it('routes anonymous log in CTA to /auth and closes menu', () => {
    const onClose = vi.fn()
    mockedUseAuth.mockReturnValue({
      status: 'anonymous',
      email: null,
      isAuthenticated: false,
      signOut: vi.fn(),
    })

    render(
      <MemoryRouter initialEntries={['/']}>
        <Routes>
          <Route
            path="/"
            element={
              <BurgerMenu
                open
                onClose={onClose}
                triggerRef={{ current: document.createElement('button') }}
              />
            }
          />
          <Route path="/auth" element={<h1>Auth Destination</h1>} />
        </Routes>
      </MemoryRouter>,
    )

    fireEvent.click(screen.getByRole('button', { name: 'Log in' }))
    expect(screen.getByRole('heading', { name: 'Auth Destination' })).toBeInTheDocument()
    expect(onClose).toHaveBeenCalledTimes(1)
  })
})
