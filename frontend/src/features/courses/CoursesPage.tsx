import { useMutation, useQueryClient } from '@tanstack/react-query'
import { FolderPlus, PenSquare, Plus, Trash2 } from 'lucide-react'
import { FormEvent, useMemo, useState } from 'react'
import {
  addBlockElement,
  addCourseBlock,
  changeElementCompletionMode,
  createCourse,
  removeBlockElement,
  removeCourseBlock,
  renameCourse,
} from '../../api/courses'
import { Button } from '../../components/ui/Button'
import { PageHeader } from '../../components/ui/PageHeader'
import { ApiError } from '../../types/api'
import { useCourses } from './useCourses'
import styles from './CoursesPage.module.css'

type ElementDraft = {
  id: string
  title: string
  type: 'test' | 'lecture_material' | 'download_file'
  completionMode: 'none' | 'manual'
}

type BlockDraft = {
  id: string
  title: string
  elements: ElementDraft[]
}

type NewElementState = {
  title: string
  type: 'test' | 'lecture_material' | 'download_file'
  completionMode: 'none' | 'manual'
  fileName: string
  sizeBytes: number
  quizId: string
}

function defaultNewElement(): NewElementState {
  return {
    title: '',
    type: 'lecture_material',
    completionMode: 'none',
    fileName: 'material.pdf',
    sizeBytes: 1024,
    quizId: '',
  }
}

export function CoursesPage() {
  const queryClient = useQueryClient()
  const { data: courses = [], error, isLoading, refetch } = useCourses()
  const [courseTitle, setCourseTitle] = useState('')
  const [newBlockTitle, setNewBlockTitle] = useState<Record<string, string>>({})
  const [newElement, setNewElement] = useState<Record<string, NewElementState>>({})
  const [courseBlocks, setCourseBlocks] = useState<Record<string, BlockDraft[]>>({})
  const [selectedCourseID, setSelectedCourseID] = useState<string | null>(null)
  const [actionError, setActionError] = useState<string | null>(null)

  const activeCourseID = selectedCourseID ?? courses[0]?.ID ?? null
  const activeCourse = useMemo(() => courses.find((course) => course.ID === activeCourseID) ?? null, [courses, activeCourseID])

  const createCourseMutation = useMutation({
    mutationFn: createCourse,
    onSuccess: () => {
      setCourseTitle('')
      void queryClient.invalidateQueries({ queryKey: ['courses'] })
    },
  })

  const renameCourseMutation = useMutation({
    mutationFn: ({ courseId, title }: { courseId: string; title: string }) => renameCourse(courseId, { title }),
    onSuccess: () => void queryClient.invalidateQueries({ queryKey: ['courses'] }),
  })

  const addBlockMutation = useMutation({
    mutationFn: ({ courseId, title }: { courseId: string; title: string }) => addCourseBlock(courseId, title),
    onSuccess: (data, variables) => {
      setCourseBlocks((prev) => {
        const next = prev[variables.courseId] ? [...prev[variables.courseId]] : []
        next.push({ id: data.id, title: variables.title, elements: [] })
        return { ...prev, [variables.courseId]: next }
      })
      setNewBlockTitle((prev) => ({ ...prev, [variables.courseId]: '' }))
      void queryClient.invalidateQueries({ queryKey: ['courses'] })
    },
  })

  const deleteBlockMutation = useMutation({
    mutationFn: ({ courseId, blockId }: { courseId: string; blockId: string }) => removeCourseBlock(courseId, blockId),
    onSuccess: (_, variables) => {
      setCourseBlocks((prev) => ({
        ...prev,
        [variables.courseId]: (prev[variables.courseId] ?? []).filter((block) => block.id !== variables.blockId),
      }))
      void queryClient.invalidateQueries({ queryKey: ['courses'] })
    },
  })

  const addElementMutation = useMutation({
    mutationFn: ({
      blockId,
      payload,
    }: {
      blockId: string
      payload: {
        title: string
        type: 'test' | 'lecture_material' | 'download_file'
        completion_mode: 'none' | 'manual'
        file_name?: string
        size_bytes?: number
        quiz_id?: string
      }
    }) => addBlockElement(blockId, payload),
    onSuccess: (data, variables) => {
      setCourseBlocks((prev) => {
        const next = { ...prev }
        for (const [courseId, blocks] of Object.entries(next)) {
          next[courseId] = blocks.map((block) =>
            block.id === variables.blockId
              ? {
                  ...block,
                  elements: [
                    ...block.elements,
                    {
                      id: data.id,
                      title: variables.payload.title,
                      type: variables.payload.type,
                      completionMode: variables.payload.completion_mode,
                    },
                  ],
                }
              : block,
          )
        }
        return next
      })
      setNewElement((prev) => ({ ...prev, [variables.blockId]: defaultNewElement() }))
    },
  })

  const deleteElementMutation = useMutation({
    mutationFn: ({ blockId, elementId }: { blockId: string; elementId: string }) => removeBlockElement(blockId, elementId),
    onSuccess: (_, variables) => {
      setCourseBlocks((prev) => {
        const next = { ...prev }
        for (const [courseId, blocks] of Object.entries(next)) {
          next[courseId] = blocks.map((block) =>
            block.id === variables.blockId
              ? { ...block, elements: block.elements.filter((element) => element.id !== variables.elementId) }
              : block,
          )
        }
        return next
      })
    },
  })

  const editElementMutation = useMutation({
    mutationFn: ({ elementId, mode }: { elementId: string; mode: 'none' | 'manual' }) =>
      changeElementCompletionMode(elementId, mode),
  })

  function extractErrorMessage(caught: unknown, fallback: string) {
    return caught instanceof ApiError ? caught.message : fallback
  }

  async function handleCreateCourse(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setActionError(null)
    try {
      await createCourseMutation.mutateAsync({ title: courseTitle })
    } catch (caught) {
      setActionError(extractErrorMessage(caught, 'Не удалось создать курс'))
    }
  }

  async function handleRenameCourse(courseId: string, currentTitle: string) {
    const nextTitle = window.prompt('Новое название курса', currentTitle)?.trim()
    if (!nextTitle || nextTitle === currentTitle) {
      return
    }
    setActionError(null)
    try {
      await renameCourseMutation.mutateAsync({ courseId, title: nextTitle })
    } catch (caught) {
      setActionError(extractErrorMessage(caught, 'Не удалось переименовать курс'))
    }
  }

  async function handleAddBlock(courseId: string) {
    const title = (newBlockTitle[courseId] ?? '').trim()
    if (!title) {
      return
    }
    setActionError(null)
    try {
      await addBlockMutation.mutateAsync({ courseId, title })
    } catch (caught) {
      setActionError(extractErrorMessage(caught, 'Не удалось добавить блок'))
    }
  }

  async function handleDeleteBlock(courseId: string, blockId: string) {
    setActionError(null)
    try {
      await deleteBlockMutation.mutateAsync({ courseId, blockId })
    } catch (caught) {
      setActionError(extractErrorMessage(caught, 'Не удалось удалить блок'))
    }
  }

  async function handleAddElement(blockId: string) {
    const state = newElement[blockId] ?? defaultNewElement()
    if (!state.title.trim()) {
      return
    }

    const payload: {
      title: string
      type: 'test' | 'lecture_material' | 'download_file'
      completion_mode: 'none' | 'manual'
      file_name?: string
      size_bytes?: number
      quiz_id?: string
    } = {
      title: state.title.trim(),
      type: state.type,
      completion_mode: state.completionMode,
    }

    if (state.type === 'test') {
      payload.quiz_id = state.quizId.trim()
    } else {
      payload.file_name = state.fileName.trim()
      payload.size_bytes = Number(state.sizeBytes)
    }

    setActionError(null)
    try {
      await addElementMutation.mutateAsync({ blockId, payload })
    } catch (caught) {
      setActionError(extractErrorMessage(caught, 'Не удалось добавить элемент'))
    }
  }

  async function handleDeleteElement(blockId: string, elementId: string) {
    setActionError(null)
    try {
      await deleteElementMutation.mutateAsync({ blockId, elementId })
    } catch (caught) {
      setActionError(extractErrorMessage(caught, 'Не удалось удалить элемент'))
    }
  }

  async function handleEditElementMode(blockId: string, element: ElementDraft) {
    const mode = element.completionMode === 'manual' ? 'none' : 'manual'
    setActionError(null)
    try {
      await editElementMutation.mutateAsync({ elementId: element.id, mode })
      setCourseBlocks((prev) => {
        const next = { ...prev }
        for (const [courseId, blocks] of Object.entries(next)) {
          next[courseId] = blocks.map((block) =>
            block.id === blockId
              ? {
                  ...block,
                  elements: block.elements.map((item) =>
                    item.id === element.id ? { ...item, completionMode: mode } : item,
                  ),
                }
              : block,
          )
        }
        return next
      })
    } catch (caught) {
      setActionError(extractErrorMessage(caught, 'Не удалось изменить режим завершения элемента'))
    }
  }

  return (
    <>
      <PageHeader title="Курсы" description="Создавайте курсы, блоки и элементы учебного контента." />

      <form className={styles.createCourse} onSubmit={handleCreateCourse}>
        <input
          value={courseTitle}
          onChange={(event) => setCourseTitle(event.target.value)}
          placeholder="Название нового курса"
          aria-label="Название курса"
          required
        />
        <Button type="submit" disabled={createCourseMutation.isPending}>
          <Plus size={16} aria-hidden="true" />
          Создать курс
        </Button>
      </form>

      {actionError ? <p className={styles.error}>{actionError}</p> : null}
      {isLoading ? <p className={styles.state}>Загрузка курсов</p> : null}
      {error ? (
        <div className={styles.state}>
          <p>{error instanceof ApiError ? error.message : 'Не удалось загрузить курсы'}</p>
          <Button variant="secondary" onClick={() => void refetch()}>
            Повторить
          </Button>
        </div>
      ) : null}
      {courses.length === 0 ? <p className={styles.state}>Курсы не найдены.</p> : null}

      {courses.length > 0 ? (
        <section className={styles.builder}>
          <aside className={styles.courseList}>
            {courses.map((course) => (
              <button
                key={course.ID}
                type="button"
                onClick={() => setSelectedCourseID(course.ID)}
                className={course.ID === activeCourseID ? `${styles.courseCard} ${styles.active}` : styles.courseCard}
              >
                <span>{course.Title}</span>
                <small>Блоков: {course.BlocksCount}</small>
              </button>
            ))}
          </aside>

          <div className={styles.workspace}>
            {activeCourse ? (
              <>
                <div className={styles.courseHead}>
                  <h2>{activeCourse.Title}</h2>
                  <Button variant="secondary" onClick={() => void handleRenameCourse(activeCourse.ID, activeCourse.Title)}>
                    <PenSquare size={16} aria-hidden="true" />
                    Редактировать
                  </Button>
                </div>

                <div className={styles.inlineRow}>
                  <input
                    value={newBlockTitle[activeCourse.ID] ?? ''}
                    onChange={(event) =>
                      setNewBlockTitle((prev) => ({
                        ...prev,
                        [activeCourse.ID]: event.target.value,
                      }))
                    }
                    placeholder="Название блока"
                    aria-label="Название блока"
                  />
                  <Button onClick={() => void handleAddBlock(activeCourse.ID)} disabled={addBlockMutation.isPending}>
                    <FolderPlus size={16} aria-hidden="true" />
                    Добавить в блок
                  </Button>
                </div>

                <div className={styles.blocks}>
                  {(courseBlocks[activeCourse.ID] ?? []).map((block) => {
                    const elementState = newElement[block.id] ?? defaultNewElement()
                    return (
                      <article key={block.id} className={styles.blockCard}>
                        <div className={styles.blockHead}>
                          <strong>{block.title}</strong>
                          <Button variant="secondary" onClick={() => void handleDeleteBlock(activeCourse.ID, block.id)}>
                            <Trash2 size={16} aria-hidden="true" />
                            Удалить
                          </Button>
                        </div>

                        <div className={styles.inlineRow}>
                          <input
                            value={elementState.title}
                            onChange={(event) =>
                              setNewElement((prev) => ({
                                ...prev,
                                [block.id]: { ...elementState, title: event.target.value },
                              }))
                            }
                            placeholder="Название элемента"
                          />
                          <select
                            value={elementState.type}
                            onChange={(event) =>
                              setNewElement((prev) => ({
                                ...prev,
                                [block.id]: {
                                  ...elementState,
                                  type: event.target.value as NewElementState['type'],
                                },
                              }))
                            }
                          >
                            <option value="lecture_material">Лекция</option>
                            <option value="download_file">Файл</option>
                            <option value="test">Тест</option>
                          </select>
                          <Button onClick={() => void handleAddElement(block.id)} disabled={addElementMutation.isPending}>
                            <Plus size={16} aria-hidden="true" />
                            Добавить
                          </Button>
                        </div>

                        {elementState.type === 'test' ? (
                          <input
                            className={styles.extraInput}
                            value={elementState.quizId}
                            onChange={(event) =>
                              setNewElement((prev) => ({
                                ...prev,
                                [block.id]: { ...elementState, quizId: event.target.value },
                              }))
                            }
                            placeholder="ID теста (quiz_id)"
                          />
                        ) : (
                          <div className={styles.inlineRow}>
                            <input
                              value={elementState.fileName}
                              onChange={(event) =>
                                setNewElement((prev) => ({
                                  ...prev,
                                  [block.id]: { ...elementState, fileName: event.target.value },
                                }))
                              }
                              placeholder="Имя файла"
                            />
                            <input
                              type="number"
                              min={1}
                              value={elementState.sizeBytes}
                              onChange={(event) =>
                                setNewElement((prev) => ({
                                  ...prev,
                                  [block.id]: { ...elementState, sizeBytes: Number(event.target.value) || 1 },
                                }))
                              }
                              placeholder="Размер в байтах"
                            />
                          </div>
                        )}

                        <ul className={styles.elements}>
                          {block.elements.map((element) => (
                            <li key={element.id} className={styles.elementRow}>
                              <div>
                                <strong>{element.title}</strong>
                                <small>
                                  {element.type} · completion: {element.completionMode}
                                </small>
                              </div>
                              <div className={styles.elementActions}>
                                <Button
                                  variant="secondary"
                                  onClick={() => void handleEditElementMode(block.id, element)}
                                  title="Редактировать"
                                >
                                  <PenSquare size={16} aria-hidden="true" />
                                </Button>
                                <Button
                                  variant="secondary"
                                  onClick={() => void handleDeleteElement(block.id, element.id)}
                                  title="Удалить"
                                >
                                  <Trash2 size={16} aria-hidden="true" />
                                </Button>
                              </div>
                            </li>
                          ))}
                        </ul>
                      </article>
                    )
                  })}
                </div>
              </>
            ) : null}
          </div>
        </section>
      ) : null}
    </>
  )
}
