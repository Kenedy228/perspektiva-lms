import { useEffect } from 'react'
import { Navigate, Outlet } from 'react-router-dom'
import { clearSessionToken } from '../../api/tokenStore'
import { useSession } from './useSession'
import styles from './StudentGuard.module.css'

export function OrgGuard() {
  const { isLoading, isAuthenticated, isUnauthenticated, isOrganization, isManager, isStudent } = useSession()

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
  if (isStudent) return <Navigate to="/my/courses" replace />

  if (!isOrganization) {
    return (
      <div className={styles.center}>
        <p className={styles.title}>Нет доступа</p>
        <p className={styles.hint}>Этот раздел доступен только для учётных записей организаций.</p>
      </div>
    )
  }

  return <Outlet />
}
