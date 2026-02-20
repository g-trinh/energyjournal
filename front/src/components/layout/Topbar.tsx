import { useEffect, useRef, useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import BurgerMenu from './BurgerMenu'
import LogoMark from './LogoMark'
import NavMenu from './NavMenu'
import UserMenu from './UserMenu'
import '@/styles/topbar.css'

export default function Topbar() {
  const location = useLocation()
  const [burgerOpen, setBurgerOpen] = useState(false)
  const triggerRef = useRef<HTMLButtonElement>(null)

  useEffect(() => {
    setBurgerOpen(false)
  }, [location.pathname])

  return (
    <header className="topbar" role="banner">
      <div className="topbar-inner">
        <div className="topbar-left">
          <Link to="/" className="topbar-brand-link" aria-label="Energy Journal home">
            <div className="topbar-brand" aria-label="Energy Journal">
              <LogoMark size={36} />
              <span className="topbar-brand-name">Energy Journal</span>
            </div>
          </Link>

          <div className="topbar-brand-divider" />
          <NavMenu />
        </div>

        <div className="topbar-desktop-user">
          <UserMenu />
        </div>

        <button
          ref={triggerRef}
          type="button"
          className="topbar-mobile-toggle"
          aria-haspopup="dialog"
          aria-expanded={burgerOpen}
          aria-label={burgerOpen ? 'Close menu' : 'Open menu'}
          onClick={() => setBurgerOpen((current) => !current)}
        >
          {burgerOpen ? '×' : '☰'}
        </button>
      </div>

      <BurgerMenu open={burgerOpen} onClose={() => setBurgerOpen(false)} triggerRef={triggerRef} />
    </header>
  )
}
