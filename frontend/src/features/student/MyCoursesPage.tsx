import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { BookOpen } from 'lucide-react'
import { Link } from 'react-router-dom'
import { getCourses } from '../../api/courses'
import { PageHeader } from '../../components/ui/PageHeader'
import { Pagination } from '../../components/ui/Pagination'
import styles from './MyCoursesPage.module.css'

const LIMIT = 20

export function MyCoursesPage() {
  const [search, setSearch] = useState('')
  const [offset, setOffset] = useState(0)

  const coursesQuery = useQuery({
    queryKey: ['my-courses', search, offset],
    queryFn: () => getCourses({ title: search || undefined, limit: LIMIT, offset }),
  })

  const courses = coursesQuery.data ?? []
  const hasMore = courses.length >= LIMIT

  return (
    <>
      <PageHeader title="Мои курсы" description="Курсы, на которые вы зачислены." />

      <div className={styles.toolbar}>
        <input
          className={styles.searchInput}
          value={search}
          onChange={(e) => { setSearch(e.target.value); setOffset(0) }}
          placeholder="Поиск по названию"
        />
      </div>

      {coursesQuery.isLoading ? (
        <p className={styles.state}>Загрузка…</p>
      ) : coursesQuery.isError ? (
        <p className={styles.stateError}>Не удалось загрузить курсы.</p>
      ) : courses.length === 0 ? (
        <div className={styles.empty}>
          <BookOpen size={40} className={styles.emptyIcon} aria-hidden="true" />
          <p>У вас пока нет доступных курсов.</p>
        </div>
      ) : (
        <>
          <section className={styles.cards}>
            {courses.map((course) => (
              <Link key={course.ID} to={`/my/courses/${course.ID}`} className={styles.cardLink}>
                <article className={styles.card}>
                  <div className={styles.cardInfo}>
                    <h3 className={styles.cardTitle}>{course.Title}</h3>
                    <p className={styles.cardMeta}>Блоков: {course.BlocksCount}</p>
                  </div>
                </article>
              </Link>
            ))}
          </section>

          {hasMore || offset > 0 ? (
            <Pagination
              offset={offset}
              limit={LIMIT}
              hasMore={hasMore}
              onOffsetChange={setOffset}
            />
          ) : null}
        </>
      )}
    </>
  )
}
