import { useQueryClient } from '@tanstack/react-query'
import { BarChart2, BookMarked, BookOpen, Building2, ClipboardList, Database, GraduationCap, LogOut, UsersRound, UserCog } from 'lucide-react'
import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { logout } from '../../api/auth'
import { clearSessionToken } from '../../api/tokenStore'
import { useSession } from '../../features/auth/useSession'
import styles from './AppLayout.module.css'

type NavItem = { to: string; label: string; icon: React.ElementType; adminOnly?: boolean }

const navItems: NavItem[] = [
  { to: '/courses', label: 'Курсы', icon: BookOpen },
  { to: '/banks', label: 'Банки вопросов', icon: Database },
  { to: '/quizzes', label: 'Тесты', icon: ClipboardList },
  { to: '/organizations', label: 'Организации', icon: Building2, adminOnly: true },
  { to: '/people', label: 'Сотрудники', icon: UsersRound, adminOnly: true },
  { to: '/accounts', label: 'Учётные записи', icon: UserCog, adminOnly: true },
  { to: '/enrollments', label: 'Зачисления', icon: BookMarked, adminOnly: true },
  { to: '/statistics', label: 'Статистика', icon: BarChart2, adminOnly: true },
]

const roleLabels: Record<string, string> = {
  admin: 'Администратор',
  creator: 'Создатель',
}

export function AppLayout() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { role, isAdmin } = useSession()

  async function handleLogout() {
    try {
      await logout()
    } catch {
      clearSessionToken()
    }
    queryClient.clear()
    navigate('/login', { replace: true })
  }

  const visibleNavItems = navItems.filter((item) => !item.adminOnly || isAdmin)

  return (
    <div className={styles.shell}>
      <aside className={styles.sidebar}>
        <div className={styles.brand}>
          <GraduationCap aria-hidden="true" />
          <span>LMS</span>
        </div>
        <nav className={styles.nav} aria-label="Primary navigation">
          {visibleNavItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.to === '/'}
              className={({ isActive }) => (isActive ? `${styles.navLink} ${styles.active}` : styles.navLink)}
            >
              <item.icon aria-hidden="true" size={18} />
              <span>{item.label}</span>
            </NavLink>
          ))}
        </nav>
        <div className={styles.sidebarFooter}>
          {role ? <span className={styles.roleLabel}>{roleLabels[role] ?? role}</span> : null}
          <button
            type="button"
            className={styles.logoutBtn}
            onClick={() => void handleLogout()}
          >
            <LogOut size={16} aria-hidden="true" />
            Выйти
          </button>
        </div>
      </aside>
      <main className={styles.main}>
        <Outlet />
      </main>
    </div>
  )
}
