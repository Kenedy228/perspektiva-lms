import { useMutation, useQueryClient } from '@tanstack/react-query'
import { FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { createBank } from '../../api/banks'
import { Button } from '../../components/ui/Button'
import { PageHeader } from '../../components/ui/PageHeader'
import { ApiError } from '../../types/api'
import styles from './CreateBankPage.module.css'

export function CreateBankPage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [title, setTitle] = useState('')
  const [error, setError] = useState<string | null>(null)

  const createMutation = useMutation({
    mutationFn: (value: string) => createBank(value),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['banks'] })
      navigate('/banks')
    },
  })

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError(null)
    try {
      await createMutation.mutateAsync(title.trim())
    } catch (caught) {
      setError(caught instanceof ApiError ? caught.message : 'Не удалось создать банк')
    }
  }

  return (
    <>
      <PageHeader title="Создание банка" description="Укажите название нового банка вопросов." />
      <form className={styles.form} onSubmit={handleSubmit}>
        <label htmlFor="title">Название банка</label>
        <input id="title" value={title} onChange={(event) => setTitle(event.target.value)} required />
        {error ? <p className={styles.error}>{error}</p> : null}
        <div className={styles.actions}>
          <Button type="submit" disabled={createMutation.isPending}>
            Создать банк
          </Button>
          <Button type="button" variant="secondary" onClick={() => navigate('/banks')}>
            Отмена
          </Button>
        </div>
      </form>
    </>
  )
}
