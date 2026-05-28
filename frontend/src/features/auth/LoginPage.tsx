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
      <section className={styles.shell}>
        <form className={styles.form} onSubmit={handleSubmit} noValidate>
          <h2>Вход в систему</h2>
          <p className={styles.subtitle}>Используйте учетную запись организации</p>

          <label htmlFor="login">Логин</label>
          <input
            id="login"
            value={loginValue}
            onChange={(event) => setLoginValue(event.target.value)}
            autoComplete="username"
            placeholder="Введите логин"
            required
          />

          <label htmlFor="password">Пароль</label>
          <input
            id="password"
            value={password}
            type="password"
            onChange={(event) => setPassword(event.target.value)}
            autoComplete="current-password"
            placeholder="Введите пароль"
            required
          />

          {error ? <p className={styles.error}>{error}</p> : null}

          <Button type="submit" disabled={isSubmitting} className={styles.submit}>
            {isSubmitting ? 'Выполняется вход' : 'Войти'}
          </Button>
        </form>
      </section>
    </main>
  )
}
