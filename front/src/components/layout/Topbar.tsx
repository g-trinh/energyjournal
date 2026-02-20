import { useRef, useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import BurgerMenu from './BurgerMenu'
import LogoMark from './LogoMark'
import NavMenu from './NavMenu'
import UserMenu from './UserMenu'
import '@/styles/topbar.css'

export default function Topbar() {
  const location = useLocation()
  const [burgerMenuState, setBurgerMenuState] = useState<{ open: boolean; path: string }>({
    open: false,
    path: location.pathname,
  })
  const triggerRef = useRef<HTMLButtonElement>(null)

  const burgerOpen = burgerMenuState.open && burgerMenuState.path === location.pathname

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
          onClick={() =>
            setBurgerMenuState((current) => ({
              open: !current.open,
              path: location.pathname,
            }))
          }
        >
          {burgerOpen ? '×' : '☰'}
        </button>
      </div>

      <BurgerMenu
        open={burgerOpen}
        onClose={() => setBurgerMenuState((current) => ({ ...current, open: false }))}
        triggerRef={triggerRef}
      />
    </header>
  )
}
