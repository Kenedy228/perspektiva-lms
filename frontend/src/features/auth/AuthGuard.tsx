import { useEffect } from 'react'
import { Navigate, Outlet } from 'react-router-dom'
import { clearSessionToken } from '../../api/tokenStore'
import { useSession } from './useSession'
import styles from './AuthGuard.module.css'

export function AuthGuard() {
  const { isLoading, isAuthenticated, isUnauthenticated, isManager, isStudent, isOrganization } = useSession()

  useEffect(() => {
    if (isUnauthenticated) {
      clearSessionToken()
    }
  }, [isUnauthenticated])

  if (isLoading) {
    return (
      <div className={styles.center}>
        <p className={styles.hint}>Загрузка…</p>
      </div>
    )
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (!isManager) {
    if (isStudent) return <Navigate to="/my/courses" replace />
    if (isOrganization) return <Navigate to="/org/statistics" replace />
    return (
      <div className={styles.center}>
        <p className={styles.title}>Нет доступа</p>
        <p className={styles.hint}>Эта панель управления доступна только администраторам и создателям.</p>
      </div>
    )
  }

  return <Outlet />
}
