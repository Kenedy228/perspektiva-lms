import { useQuery } from '@tanstack/react-query'
import { ChevronRight, Plus } from 'lucide-react'
import { Link } from 'react-router-dom'
import { listBanks } from '../../api/banks'
import { Button } from '../../components/ui/Button'
import { PageHeader } from '../../components/ui/PageHeader'
import styles from './BanksPage.module.css'
import { useState } from 'react'

const BANK_PAGE_SIZE = 8

export function BanksPage() {
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)
  const banksQuery = useQuery({
    queryKey: ['banks', search],
    queryFn: () => listBanks({ title: search, limit: 200, offset: 0 }),
  })

  const allBanks = banksQuery.data ?? []
  const totalPages = Math.max(1, Math.ceil(allBanks.length / BANK_PAGE_SIZE))
  const safePage = Math.min(page, totalPages)
  const banks = allBanks.slice((safePage - 1) * BANK_PAGE_SIZE, safePage * BANK_PAGE_SIZE)

  return (
    <>
      <PageHeader title="Банки вопросов" description="Список всех банков вопросов." />

      <div className={styles.toolbar}>
        <input
          value={search}
          onChange={(event) => {
            setSearch(event.target.value)
            setPage(1)
          }}
          placeholder="Поиск по названию банка"
        />
        <Button as={Link} to="/banks/new">
          <Plus size={16} />
          Добавить
        </Button>
      </div>

      {allBanks.length === 0 ? <p className={styles.empty}>Банков вопросов пока нет.</p> : null}

      {allBanks.length > 0 ? (
        <>
          <section className={styles.cards}>
            {banks.map((bank) => (
              <Link key={bank.ID} to={`/banks/${bank.ID}`} className={styles.card}>
                <div className={styles.cardInfo}>
                  <h3>{bank.Title}</h3>
                  <p>Вопросов: {bank.QuestionsCount}</p>
                </div>
                <ChevronRight size={18} className={styles.cardArrow} aria-hidden="true" />
              </Link>
            ))}
          </section>

          <div className={styles.pager}>
            <Button variant="secondary" disabled={safePage <= 1} onClick={() => setPage((prev) => Math.max(1, prev - 1))}>
              Назад
            </Button>
            <span>
              {safePage}/{totalPages}
            </span>
            <Button
              variant="secondary"
              disabled={safePage >= totalPages}
              onClick={() => setPage((prev) => Math.min(totalPages, prev + 1))}
            >
              Вперед
            </Button>
          </div>
        </>
      ) : null}
    </>
  )
}
