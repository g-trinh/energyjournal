import { fireEvent, render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import NavMenu from './NavMenu'
import { useAuth } from '@/contexts/AuthContext'

vi.mock('@/contexts/AuthContext', () => ({
  useAuth: vi.fn(),
}))

describe('NavMenu', () => {
  const mockedUseAuth = vi.mocked(useAuth)

  beforeEach(() => {
    mockedUseAuth.mockReset()
  })

  it('shows active state on /timespending for authenticated users', () => {
    mockedUseAuth.mockReturnValue({
      status: 'authenticated',
      email: 'person@example.com',
      isAuthenticated: true,
      signOut: vi.fn(),
    })

    render(
      <MemoryRouter initialEntries={['/timespending']}>
        <NavMenu />
      </MemoryRouter>,
    )

    const item = screen.getByRole('button', { name: 'Time Spending' })
    expect(item).toHaveClass('topbar-nav-item-active')
    expect(item).toHaveAttribute('aria-current', 'page')
  })

  it('redirects anonymous users to /auth when clicking Time Spending', () => {
    mockedUseAuth.mockReturnValue({
      status: 'anonymous',
      email: null,
      isAuthenticated: false,
      signOut: vi.fn(),
    })

    render(
      <MemoryRouter initialEntries={['/']}>
        <Routes>
          <Route path="/" element={<NavMenu />} />
          <Route path="/auth" element={<h1>Auth Destination</h1>} />
        </Routes>
      </MemoryRouter>,
    )

    fireEvent.click(screen.getByRole('button', { name: 'Time Spending' }))
    expect(screen.getByRole('heading', { name: 'Auth Destination' })).toBeInTheDocument()
  })
})
