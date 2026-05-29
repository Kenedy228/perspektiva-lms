import { ChevronLeft, ChevronRight } from 'lucide-react'
import styles from './Pagination.module.css'

type PaginationProps = {
  offset: number
  limit: number
  onOffsetChange: (offset: number) => void
} & (
  | { total: number; hasMore?: never }
  | { total?: never; hasMore: boolean }
)

export function Pagination({ offset, limit, onOffsetChange, total, hasMore }: PaginationProps) {
  const hasPrev = offset > 0
  const hasNext = total !== undefined ? offset + limit < total : hasMore

  const from = (total === 0) ? 0 : offset + 1
  const to = total !== undefined ? Math.min(offset + limit, total) : offset + limit

  return (
    <div className={styles.pagination}>
      <button
        type="button"
        className={styles.btn}
        onClick={() => onOffsetChange(Math.max(0, offset - limit))}
        disabled={!hasPrev}
        aria-label="Предыдущая страница"
      >
        <ChevronLeft size={16} aria-hidden="true" />
      </button>

      <span className={styles.info}>
        {total === 0
          ? 'Нет результатов'
          : total !== undefined
            ? `${from}–${to} из ${total}`
            : `${from}–${to}`}
      </span>

      <button
        type="button"
        className={styles.btn}
        onClick={() => onOffsetChange(offset + limit)}
        disabled={!hasNext}
        aria-label="Следующая страница"
      >
        <ChevronRight size={16} aria-hidden="true" />
      </button>
    </div>
  )
}
