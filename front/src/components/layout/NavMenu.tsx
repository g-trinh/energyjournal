import { useAuth } from '@/contexts/AuthContext'
import { cn } from '@/lib/utils'
import { useLocation, useNavigate } from 'react-router-dom'

export default function NavMenu() {
  const navigate = useNavigate()
  const location = useLocation()
  const { status } = useAuth()

  const isAnonymous = status === 'anonymous'
  const isActive = !isAnonymous && location.pathname === '/timespending'

  function handleTimeSpendingClick() {
    navigate(isAnonymous ? '/auth' : '/timespending')
  }

  return (
    <nav aria-label="Main navigation" className="topbar-nav">
      <ul role="list" className="topbar-nav-list">
        <li>
          <button
            type="button"
            className={cn(
              'topbar-nav-item',
              isActive && 'topbar-nav-item-active',
              isAnonymous && 'topbar-nav-item-anonymous',
            )}
            aria-current={isActive ? 'page' : undefined}
            onClick={handleTimeSpendingClick}
          >
            Time Spending
          </button>
        </li>
      </ul>
    </nav>
  )
}
