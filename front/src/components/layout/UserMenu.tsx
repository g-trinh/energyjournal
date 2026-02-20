import { useAuth } from '@/contexts/AuthContext'
import { useEffect, useRef, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'

export default function UserMenu() {
  const navigate = useNavigate()
  const location = useLocation()
  const { status, email, signOut } = useAuth()
  const [open, setOpen] = useState(false)
  const rootRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    setOpen(false)
  }, [location.pathname])

  useEffect(() => {
    if (!open) {
      return
    }

    function onKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        setOpen(false)
      }
    }

    function onPointerDown(event: MouseEvent) {
      if (!rootRef.current?.contains(event.target as Node)) {
        setOpen(false)
      }
    }

    document.addEventListener('keydown', onKeyDown)
    document.addEventListener('mousedown', onPointerDown)

    return () => {
      document.removeEventListener('keydown', onKeyDown)
      document.removeEventListener('mousedown', onPointerDown)
    }
  }, [open])

  if (status === 'loading') {
    return <div aria-label="Loading user menu" className="topbar-user-skeleton" />
  }

  if (status === 'anonymous') {
    return (
      <button
        type="button"
        className="topbar-login-btn"
        onClick={() => navigate('/auth')}
      >
        Log in
      </button>
    )
  }

  return (
    <div className="topbar-user" ref={rootRef}>
      <button
        type="button"
        className="topbar-user-trigger"
        aria-haspopup="menu"
        aria-expanded={open}
        onClick={() => setOpen((current) => !current)}
      >
        <span className="topbar-user-trigger-email">{email ?? 'Account'}</span>
      </button>

      {open && (
        <div className="topbar-user-menu" role="menu" aria-label="Account actions">
          <p className="topbar-user-menu-email">{email ?? 'Account'}</p>
          <div className="topbar-user-menu-divider" />
          <button
            type="button"
            className="topbar-user-menu-signout"
            role="menuitem"
            onClick={signOut}
          >
            Sign out
          </button>
        </div>
      )}
    </div>
  )
}
