import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import LandingPage from './LandingPage'

describe('LandingPage', () => {
  it('renders key sections and chart preview', () => {
    render(
      <MemoryRouter>
        <LandingPage />
      </MemoryRouter>,
    )

    expect(screen.getByRole('heading', { name: 'Three dimensions of energy' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: 'Simple to start, powerful over time' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: 'Start understanding your energy today' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: /7-day energy preview/i })).toBeInTheDocument()
  })

  it('routes both landing ctas to /timespending', () => {
    render(
      <MemoryRouter>
        <LandingPage />
      </MemoryRouter>,
    )

    expect(screen.getByRole('link', { name: 'Start tracking →' })).toHaveAttribute('href', '/timespending')
    expect(screen.getByRole('link', { name: 'Create free account →' })).toHaveAttribute('href', '/timespending')
  })
})
