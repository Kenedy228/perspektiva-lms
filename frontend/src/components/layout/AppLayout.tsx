import { BookOpen, Building2, Database, GraduationCap, UsersRound, UserCog } from 'lucide-react'
import { NavLink, Outlet } from 'react-router-dom'
import styles from './AppLayout.module.css'

const navItems = [
  { to: '/courses', label: 'Курсы', icon: BookOpen },
  { to: '/banks', label: 'Банки вопросов', icon: Database },
  { to: '/organizations', label: 'Организации', icon: Building2 },
  { to: '/people', label: 'Сотрудники', icon: UsersRound },
  { to: '/accounts', label: 'Учетные записи', icon: UserCog },
]

export function AppLayout() {
  return (
    <div className={styles.shell}>
      <aside className={styles.sidebar}>
        <div className={styles.brand}>
          <GraduationCap aria-hidden="true" />
          <span>LMS</span>
        </div>
        <nav className={styles.nav} aria-label="Primary navigation">
          {navItems.map((item) => (
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
      </aside>
      <main className={styles.main}>
        <Outlet />
      </main>
    </div>
  )
}
