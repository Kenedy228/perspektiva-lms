import { useEffect } from 'react'
import { Navigate, Outlet } from 'react-router-dom'
import { clearSessionToken } from '../../api/tokenStore'
import { useSession } from './useSession'
import styles from './StudentGuard.module.css'

export function StudentGuard() {
  const { isLoading, isAuthenticated, isUnauthenticated, isStudent, isManager } = useSession()

  useEffect(() => {
    if (isUnauthenticated) clearSessionToken()
  }, [isUnauthenticated])

  if (isLoading) {
    return (
      <div className={styles.center}>
        <p className={styles.hint}>Загрузка…</p>
      </div>
    )
  }

  if (!isAuthenticated) return <Navigate to="/login" replace />

  if (isManager) return <Navigate to="/courses" replace />

  if (!isStudent) {
    return (
      <div className={styles.center}>
        <p className={styles.title}>Нет доступа</p>
        <p className={styles.hint}>Этот раздел доступен только студентам.</p>
      </div>
    )
  }

  return <Outlet />
}
