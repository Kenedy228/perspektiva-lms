import { Navigate, Outlet } from 'react-router-dom'
import { useSession } from './useSession'
import styles from './AdminOnlyGuard.module.css'

export function AdminOnlyGuard() {
  const { isLoading, isAdmin, isManager } = useSession()

  if (isLoading) {
    return null
  }

  // Creator is authenticated and passed AuthGuard, but this section is admin-only
  if (isManager && !isAdmin) {
    return (
      <div className={styles.denied}>
        <p className={styles.title}>Нет доступа</p>
        <p className={styles.hint}>Этот раздел доступен только администраторам.</p>
      </div>
    )
  }

  // Unauthenticated — let AuthGuard handle it
  if (!isManager) {
    return <Navigate to="/" replace />
  }

  return <Outlet />
}
