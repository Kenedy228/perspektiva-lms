import { LogIn } from 'lucide-react'
import { FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { login } from '../../api/auth'
import { Button } from '../../components/ui/Button'
import { ApiError } from '../../types/api'
import styles from './LoginPage.module.css'

export function LoginPage() {
  const navigate = useNavigate()
  const [loginValue, setLoginValue] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError(null)
    setIsSubmitting(true)

    try {
      await login({ login: loginValue, password })
      navigate('/')
    } catch (caught) {
      setError(caught instanceof ApiError ? caught.message : 'Не удалось выполнить вход')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <main className={styles.page}>
      <section className={styles.panel}>
        <div className={styles.heading}>
          <LogIn aria-hidden="true" />
          <div>
            <h1>Вход в LMS</h1>
            <p>Используйте учетную запись организации.</p>
          </div>
        </div>
        <form className={styles.form} onSubmit={handleSubmit}>
          <label>
            Логин
            <input value={loginValue} onChange={(event) => setLoginValue(event.target.value)} autoComplete="username" />
          </label>
          <label>
            Пароль
            <input
              value={password}
              type="password"
              onChange={(event) => setPassword(event.target.value)}
              autoComplete="current-password"
            />
          </label>
          {error ? <p className={styles.error}>{error}</p> : null}
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Выполняется вход' : 'Войти'}
          </Button>
        </form>
      </section>
    </main>
  )
}
