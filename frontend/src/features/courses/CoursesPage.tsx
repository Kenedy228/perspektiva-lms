import { Plus } from 'lucide-react'
import { Button } from '../../components/ui/Button'
import { PageHeader } from '../../components/ui/PageHeader'
import { ApiError } from '../../types/api'
import { useCourses } from './useCourses'
import styles from './CoursesPage.module.css'

export function CoursesPage() {
  const { data: courses, error, isLoading, refetch } = useCourses()

  return (
    <>
      <PageHeader
        title="Курсы"
        description="Опубликованные версии неизменяемы; изменения выполняются через новые версии."
        actions={
          <Button>
            <Plus size={16} aria-hidden="true" />
            Новый курс
          </Button>
        }
      />
      {isLoading ? <p className={styles.state}>Загрузка курсов</p> : null}
      {error ? (
        <div className={styles.state}>
          <p>{error instanceof ApiError ? error.message : 'Не удалось загрузить курсы'}</p>
          <Button variant="secondary" onClick={() => void refetch()}>
            Повторить
          </Button>
        </div>
      ) : null}
      {courses && courses.length === 0 ? <p className={styles.state}>Курсы не найдены.</p> : null}
      {courses && courses.length > 0 ? (
        <div className={styles.tableWrap}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th>Название</th>
                <th>Опубликован</th>
                <th>Версии</th>
              </tr>
            </thead>
            <tbody>
              {courses.map((course) => (
                <tr key={course.ID}>
                  <td>{course.Title}</td>
                  <td>{course.Published ? 'Да' : 'Нет'}</td>
                  <td>{course.VersionsCount}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : null}
    </>
  )
}
