import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import Topbar from './Topbar'

vi.mock('./NavMenu', () => ({
  default: () => <nav aria-label="Primary navigation" />,
}))

vi.mock('./UserMenu', () => ({
  default: () => <div>User menu</div>,
}))

vi.mock('./BurgerMenu', () => ({
  default: () => null,
}))

describe('Topbar', () => {
  it('renders the brand as a home link', () => {
    render(
      <MemoryRouter>
        <Topbar />
      </MemoryRouter>,
    )

    const brandLink = screen.getByRole('link', { name: 'Energy Journal home' })

    expect(brandLink).toHaveAttribute('href', '/')
    expect(screen.getByText('Energy Journal')).toBeInTheDocument()
  })
})
