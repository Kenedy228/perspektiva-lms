import { useQueryClient } from '@tanstack/react-query'
import { BookMarked, BookOpen, GraduationCap, LogOut } from 'lucide-react'
import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { logout } from '../../api/auth'
import { clearSessionToken } from '../../api/tokenStore'
import styles from './AppLayout.module.css'

export function StudentLayout() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  async function handleLogout() {
    try {
      await logout()
    } catch {
      clearSessionToken()
    }
    queryClient.clear()
    navigate('/login', { replace: true })
  }

  return (
    <div className={styles.shell}>
      <aside className={styles.sidebar}>
        <div className={styles.brand}>
          <GraduationCap aria-hidden="true" />
          <span>LMS</span>
        </div>
        <nav className={styles.nav} aria-label="Primary navigation">
          <NavLink
            to="/my/courses"
            className={({ isActive }) => (isActive ? `${styles.navLink} ${styles.active}` : styles.navLink)}
          >
            <BookOpen aria-hidden="true" size={18} />
            <span>Мои курсы</span>
          </NavLink>
          <NavLink
            to="/my/enrollments"
            className={({ isActive }) => (isActive ? `${styles.navLink} ${styles.active}` : styles.navLink)}
          >
            <BookMarked aria-hidden="true" size={18} />
            <span>Мои зачисления</span>
          </NavLink>
        </nav>
        <div className={styles.sidebarFooter}>
          <span className={styles.roleLabel}>Студент</span>
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
