import { useAuth } from '@/contexts/AuthContext'
import { cn } from '@/lib/utils'
import { useEffect, useRef, type RefObject } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'

interface BurgerMenuProps {
  open: boolean
  onClose: () => void
  triggerRef: RefObject<HTMLButtonElement | null>
}

function ClockIcon({ active }: { active: boolean }) {
  return (
    <svg className={cn('topbar-mobile-row-icon', active && 'topbar-mobile-row-icon-active')} viewBox="0 0 20 20" aria-hidden="true">
      <circle cx="10" cy="10" r="7" fill="none" stroke="currentColor" strokeWidth="1.6" />
      <path d="M10 6v4l3 2" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

function EnergyIcon({ active }: { active: boolean }) {
  return (
    <svg className={cn('topbar-mobile-row-icon', active && 'topbar-mobile-row-icon-active')} viewBox="0 0 20 20" aria-hidden="true">
      <path d="M10 2.8 6 10h3l-1 7.2L14 10h-3l1.2-7.2Z" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinejoin="round" />
    </svg>
  )
}

function ChevronIcon() {
  return (
    <svg className="topbar-mobile-row-chevron" viewBox="0 0 20 20" aria-hidden="true">
      <path d="M7 4l6 6-6 6" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export default function BurgerMenu({ open, onClose, triggerRef }: BurgerMenuProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const { status, email, signOut } = useAuth()
  const dialogRef = useRef<HTMLDivElement>(null)

  const isAuthenticated = status === 'authenticated'
  const isAnonymous = status === 'anonymous'
  const isTimeSpendingActive = location.pathname === '/timespending' && !isAnonymous
  const isEnergyActive = location.pathname === '/energy/levels/edit' && !isAnonymous

  useEffect(() => {
    if (!open) {
      return
    }

    const focusable = dialogRef.current?.querySelectorAll<HTMLElement>(
      'button:not(:disabled), a[href], [tabindex]:not([tabindex="-1"])',
    )
    focusable?.[0]?.focus()

    function onKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        onClose()
        return
      }

      if (event.key !== 'Tab' || !focusable || focusable.length === 0) {
        return
      }

      const first = focusable[0]
      const last = focusable[focusable.length - 1]
      const active = document.activeElement

      if (event.shiftKey && active === first) {
        event.preventDefault()
        last.focus()
      } else if (!event.shiftKey && active === last) {
        event.preventDefault()
        first.focus()
      }
    }

    document.addEventListener('keydown', onKeyDown)
    const triggerElement = triggerRef.current

    return () => {
      document.removeEventListener('keydown', onKeyDown)
      triggerElement?.focus()
    }
  }, [open, onClose, triggerRef])

  if (!open) {
    return null
  }

  function navigateTimeSpending() {
    navigate(isAnonymous ? '/auth' : '/timespending')
    onClose()
  }

  function navigateEnergyLevels() {
    navigate(isAnonymous ? '/auth' : '/energy/levels/edit')
    onClose()
  }

  function navigateToAuth() {
    navigate('/auth')
    onClose()
  }

  function handleSignOut() {
    signOut()
    onClose()
  }

  return (
    <div
      className="topbar-mobile-overlay"
      role="dialog"
      aria-modal="true"
      aria-label="Navigation menu"
      onMouseDown={(event) => {
        if (event.target === event.currentTarget) {
          onClose()
        }
      }}
    >
      <div className="topbar-mobile-panel" ref={dialogRef}>
        {isAuthenticated && (
          <>
            <section className="topbar-mobile-user">
              <p className="topbar-mobile-user-label">SIGNED IN AS</p>
              <p className="topbar-mobile-user-email">{email ?? 'Account'}</p>
            </section>
            <div className="topbar-mobile-divider" />
          </>
        )}

        <p className="topbar-mobile-section-title">NAVIGATION</p>
        <button
          type="button"
          className={cn('topbar-mobile-row', isTimeSpendingActive && 'topbar-mobile-row-active')}
          onClick={navigateTimeSpending}
        >
          <ClockIcon active={isTimeSpendingActive} />
          <span className={cn('topbar-mobile-row-label', isAnonymous && 'topbar-mobile-row-label-anonymous')}>
            Time Spending
          </span>
          <ChevronIcon />
        </button>
        <button
          type="button"
          className={cn('topbar-mobile-row', isEnergyActive && 'topbar-mobile-row-active')}
          onClick={navigateEnergyLevels}
        >
          <EnergyIcon active={isEnergyActive} />
          <span className={cn('topbar-mobile-row-label', isAnonymous && 'topbar-mobile-row-label-anonymous')}>
            Energy Levels
          </span>
          <ChevronIcon />
        </button>

        {isAuthenticated && (
          <>
            <div className="topbar-mobile-divider" />
            <p className="topbar-mobile-section-title">ACCOUNT</p>
            <button type="button" className="topbar-mobile-row topbar-mobile-signout" onClick={handleSignOut}>
              <svg className="topbar-mobile-row-icon" viewBox="0 0 20 20" aria-hidden="true">
                <path d="M8 4H5a1 1 0 00-1 1v10a1 1 0 001 1h3" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" />
                <path d="M12 7l3 3-3 3M15 10H8" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
              <span className="topbar-mobile-row-label">Sign out</span>
            </button>
          </>
        )}

        {isAnonymous && (
          <>
            <div className="topbar-mobile-divider" />
            <button type="button" className="topbar-mobile-login-btn" onClick={navigateToAuth}>
              Log in
            </button>
          </>
        )}
      </div>
    </div>
  )
}
