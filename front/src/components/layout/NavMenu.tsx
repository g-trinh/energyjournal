import { useAuth } from '@/contexts/AuthContext'
import { cn } from '@/lib/utils'
import { useLocation, useNavigate } from 'react-router-dom'

export default function NavMenu() {
  const navigate = useNavigate()
  const location = useLocation()
  const { status } = useAuth()

  const isAnonymous = status === 'anonymous'
  const isTimeSpendingActive = !isAnonymous && location.pathname === '/timespending'
  const isEnergyActive = !isAnonymous && location.pathname === '/energy/levels/edit'

  function handleTimeSpendingClick() {
    navigate(isAnonymous ? '/auth' : '/timespending')
  }

  function handleEnergyClick() {
    navigate(isAnonymous ? '/auth' : '/energy/levels/edit')
  }

  return (
    <nav aria-label="Main navigation" className="topbar-nav">
      <ul role="list" className="topbar-nav-list">
        <li>
          <button
            type="button"
            className={cn(
              'topbar-nav-item',
              isTimeSpendingActive && 'topbar-nav-item-active',
              isAnonymous && 'topbar-nav-item-anonymous',
            )}
            aria-current={isTimeSpendingActive ? 'page' : undefined}
            onClick={handleTimeSpendingClick}
          >
            Time Spending
          </button>
        </li>
        <li>
          <button
            type="button"
            className={cn(
              'topbar-nav-item',
              isEnergyActive && 'topbar-nav-item-active',
              isAnonymous && 'topbar-nav-item-anonymous',
            )}
            aria-current={isEnergyActive ? 'page' : undefined}
            onClick={handleEnergyClick}
          >
            Energy Levels
          </button>
        </li>
      </ul>
    </nav>
  )
}
