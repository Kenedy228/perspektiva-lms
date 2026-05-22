import { Link } from 'react-router-dom'
import { Button } from '../ui/Button'
import styles from './RouteError.module.css'

export function RouteError() {
  return (
    <main className={styles.page}>
      <section className={styles.panel}>
        <h1>Страница недоступна</h1>
        <p>Не удалось открыть запрошенный раздел LMS.</p>
        <Button as={Link} to="/">
          Вернуться к обзору
        </Button>
      </section>
    </main>
  )
}
