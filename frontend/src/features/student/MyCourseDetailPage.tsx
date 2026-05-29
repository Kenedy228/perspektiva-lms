import { useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { useMutation, useQuery } from '@tanstack/react-query'
import { ArrowLeft, BookOpen, Download, FileText, PlayCircle } from 'lucide-react'
import { getCourse } from '../../api/courses'
import { listEnrollments } from '../../api/enrollments'
import { startAttempt } from '../../api/attempts'
import { Button } from '../../components/ui/Button'
import { PageHeader } from '../../components/ui/PageHeader'
import { useSession } from '../auth/useSession'
import { ApiError } from '../../types/api'
import type { CourseElementView } from '../../types/courses'
import styles from './MyCourseDetailPage.module.css'

function elementIcon(type: string) {
  if (type === 'quiz') return PlayCircle
  if (type === 'video') return BookOpen
  return type === 'document' ? Download : FileText
}

function elementTypeLabel(type: string) {
  if (type === 'quiz') return 'Тест'
  if (type === 'video') return 'Видео'
  if (type === 'document') return 'Файл'
  return 'Материал'
}

type StartState = {
  loading: boolean
  error: string | null
}

function QuizStartButton({
  element,
  enrollmentId,
  accountId,
}: {
  element: CourseElementView
  enrollmentId: string
  accountId: string
}) {
  const navigate = useNavigate()
  const [state, setState] = useState<StartState>({ loading: false, error: null })

  const startMutation = useMutation({
    mutationFn: () =>
      startAttempt({
        account_id: accountId,
        enrollment_id: enrollmentId,
        quiz_id: element.QuizID,
      }),
    onSuccess: (data) => {
      navigate(`/my/attempts/${data.id}`)
    },
    onError: (err) => {
      setState({
        loading: false,
        error: err instanceof ApiError ? err.message : 'Не удалось начать попытку',
      })
    },
  })

  return (
    <div className={styles.quizActions}>
      <Button
        onClick={() => { setState({ loading: true, error: null }); void startMutation.mutateAsync() }}
        disabled={startMutation.isPending}
      >
        <PlayCircle size={15} aria-hidden="true" />
        {startMutation.isPending ? 'Запуск…' : 'Начать тест'}
      </Button>
      {state.error ? <p className={styles.quizError}>{state.error}</p> : null}
    </div>
  )
}

export function MyCourseDetailPage() {
  const { id = '' } = useParams<{ id: string }>()
  const { accountId } = useSession()

  const courseQuery = useQuery({
    queryKey: ['course', id],
    queryFn: () => getCourse(id),
    enabled: Boolean(id),
  })

  const enrollmentQuery = useQuery({
    queryKey: ['my-enrollment-for-course', id],
    queryFn: () => listEnrollments({ course_id: id, limit: 1 }),
    enabled: Boolean(id),
  })

  const course = courseQuery.data
  const enrollment = enrollmentQuery.data?.[0]

  if (courseQuery.isLoading || enrollmentQuery.isLoading) {
    return <p className={styles.state}>Загрузка…</p>
  }

  if (courseQuery.isError || !course) {
    return (
      <div className={styles.state}>
        <p className={styles.stateError}>
          {courseQuery.error instanceof ApiError
            ? courseQuery.error.message
            : 'Курс не найден'}
        </p>
        <Button as={Link} to="/my/courses" variant="secondary">
          <ArrowLeft size={16} aria-hidden="true" />
          К курсам
        </Button>
      </div>
    )
  }

  return (
    <>
      <div className={styles.backRow}>
        <Button as={Link} to="/my/courses" variant="secondary">
          <ArrowLeft size={16} aria-hidden="true" />
          Мои курсы
        </Button>
      </div>

      <PageHeader title={course.Title} />

      {!course.Blocks || course.Blocks.length === 0 ? (
        <p className={styles.empty}>В этом курсе пока нет содержимого.</p>
      ) : (
        <div className={styles.blocks}>
          {course.Blocks.map((block, bi) => (
            <section key={block.ID} className={styles.block}>
              <h2 className={styles.blockTitle}>
                <span className={styles.blockNum}>{bi + 1}</span>
                {block.Title}
              </h2>

              {!block.Elements || block.Elements.length === 0 ? (
                <p className={styles.emptyBlock}>Блок пуст.</p>
              ) : (
                <ul className={styles.elements}>
                  {block.Elements.map((el) => {
                    const Icon = elementIcon(el.Type)
                    return (
                      <li key={el.ID} className={styles.element}>
                        <div className={styles.elementLeft}>
                          <Icon size={18} className={styles.elementIcon} aria-hidden="true" />
                          <div>
                            <p className={styles.elementTitle}>{el.Title}</p>
                            <p className={styles.elementType}>{elementTypeLabel(el.Type)}</p>
                          </div>
                        </div>
                        {el.Type === 'quiz' && el.QuizID && enrollment && accountId ? (
                          <QuizStartButton
                            element={el}
                            enrollmentId={enrollment.id}
                            accountId={accountId}
                          />
                        ) : null}
                      </li>
                    )
                  })}
                </ul>
              )}
            </section>
          ))}
        </div>
      )}
    </>
  )
}
