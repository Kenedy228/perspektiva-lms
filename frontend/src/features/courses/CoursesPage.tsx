import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  ArrowDown,
  ArrowUp,
  BarChart2,
  Download,
  FolderPlus,
  Pencil,
  Plus,
  Trash2,
  Upload,
} from 'lucide-react'
import { FormEvent, useRef, useState } from 'react'
import {
  addBlockElement,
  addCourseBlock,
  changeElementCompletionMode,
  createCourse,
  getElementDownloadURL,
  listCourseRatings,
  moveBlockElement,
  moveCourseBlock,
  removeBlockElement,
  removeCourseBlock,
  renameCourse,
  uploadElementContent,
} from '../../api/courses'
import { listBanks } from '../../api/banks'
import { createQuiz } from '../../api/quizzes'
import { Button } from '../../components/ui/Button'
import { ConfirmDialog } from '../../components/ui/ConfirmDialog'
import { FormField } from '../../components/ui/FormField'
import { Modal } from '../../components/ui/Modal'
import { PageHeader } from '../../components/ui/PageHeader'
import { ApiError } from '../../types/api'
import type { CompletionMode, ElementType } from '../../types/courses'
import type { CreateQuizPayload } from '../../types/quizzes'
import { useCourses } from './useCourses'
import styles from './CoursesPage.module.css'

// ── Local-state types ─────────────────────────────────────────────────────────

type ElementDraft = {
  id: string
  title: string
  type: ElementType
  completionMode: CompletionMode
}

type BlockDraft = {
  id: string
  courseId: string
  title: string
  elements: ElementDraft[]
}

type NewElementState = {
  title: string
  type: ElementType
  completionMode: CompletionMode
  fileName: string
  sizeBytes: number
  quizId: string
}

function defaultNewElement(): NewElementState {
  return { title: '', type: 'lecture_material', completionMode: 'none', fileName: '', sizeBytes: 1024, quizId: '' }
}

// ── CoursesPage ───────────────────────────────────────────────────────────────

export function CoursesPage() {
  const queryClient = useQueryClient()
  const { data: courses = [], isLoading, error } = useCourses()

  const [selectedCourseId, setSelectedCourseId] = useState<string | null>(null)
  const [blocks, setBlocks] = useState<BlockDraft[]>([])
  const [newBlockTitle, setNewBlockTitle] = useState('')
  const [newElement, setNewElement] = useState<Record<string, NewElementState>>({})

  const [createTitle, setCreateTitle] = useState('')
  const [pageError, setPageError] = useState<string | null>(null)

  // Rename modal
  const [renameTarget, setRenameTarget] = useState<{ id: string; title: string } | null>(null)
  const [renameError, setRenameError] = useState<string | null>(null)

  // Delete block confirm
  const [deleteBlockTarget, setDeleteBlockTarget] = useState<{ courseId: string; blockId: string; title: string } | null>(null)

  // Course ratings modal
  const [ratingsOpen, setRatingsOpen] = useState(false)

  const activeCourse = courses.find((c) => c.ID === (selectedCourseId ?? courses[0]?.ID)) ?? null
  const activeCourseId = activeCourse?.ID ?? null

  function invalidateCourses() {
    void queryClient.invalidateQueries({ queryKey: ['courses'] })
  }

  function extractMsg(caught: unknown, fallback: string) {
    return caught instanceof ApiError ? caught.message : fallback
  }

  // ── Mutations ────────────────────────────────────────────────────────────────

  const createMutation = useMutation({
    mutationFn: createCourse,
    onSuccess: () => { setCreateTitle(''); invalidateCourses() },
  })

  const renameMutation = useMutation({
    mutationFn: ({ id, title }: { id: string; title: string }) => renameCourse(id, title),
    onSuccess: () => { setRenameTarget(null); setRenameError(null); invalidateCourses() },
  })

  const addBlockMutation = useMutation({
    mutationFn: ({ courseId, title }: { courseId: string; title: string }) =>
      addCourseBlock(courseId, title),
    onSuccess: (data, vars) => {
      setBlocks((prev) => [...prev, { id: data.id, courseId: vars.courseId, title: vars.title, elements: [] }])
      setNewBlockTitle('')
      invalidateCourses()
    },
  })

  const deleteBlockMutation = useMutation({
    mutationFn: ({ courseId, blockId }: { courseId: string; blockId: string }) =>
      removeCourseBlock(courseId, blockId),
    onSuccess: (_, vars) => {
      setBlocks((prev) => prev.filter((b) => b.id !== vars.blockId))
      setDeleteBlockTarget(null)
      invalidateCourses()
    },
  })

  const moveBlockMutation = useMutation({
    mutationFn: ({ courseId, from, to }: { courseId: string; from: number; to: number }) =>
      moveCourseBlock(courseId, from, to),
    onSuccess: (_, vars) => {
      setBlocks((prev) => {
        const next = prev.filter((b) => b.courseId === vars.courseId)
        const others = prev.filter((b) => b.courseId !== vars.courseId)
        const moved = [...next]
        const [item] = moved.splice(vars.from, 1)
        moved.splice(vars.to, 0, item)
        return [...others, ...moved]
      })
    },
  })

  const addElementMutation = useMutation({
    mutationFn: ({ blockId, payload }: Parameters<typeof addBlockElement>) =>
      addBlockElement(blockId, payload),
    onSuccess: (data, vars) => {
      const state = newElement[vars.blockId] ?? defaultNewElement()
      setBlocks((prev) =>
        prev.map((b) =>
          b.id !== vars.blockId
            ? b
            : {
                ...b,
                elements: [
                  ...b.elements,
                  { id: data.id, title: state.title, type: state.type, completionMode: state.completionMode },
                ],
              },
        ),
      )
      setNewElement((prev) => ({ ...prev, [vars.blockId]: defaultNewElement() }))
    },
  })

  const deleteElementMutation = useMutation({
    mutationFn: ({ blockId, elementId }: { blockId: string; elementId: string }) =>
      removeBlockElement(blockId, elementId),
    onSuccess: (_, vars) => {
      setBlocks((prev) =>
        prev.map((b) =>
          b.id !== vars.blockId
            ? b
            : { ...b, elements: b.elements.filter((e) => e.id !== vars.elementId) },
        ),
      )
    },
  })

  const moveElementMutation = useMutation({
    mutationFn: ({ blockId, from, to }: { blockId: string; from: number; to: number }) =>
      moveBlockElement(blockId, from, to),
    onSuccess: (_, vars) => {
      setBlocks((prev) =>
        prev.map((b) => {
          if (b.id !== vars.blockId) return b
          const elems = [...b.elements]
          const [item] = elems.splice(vars.from, 1)
          elems.splice(vars.to, 0, item)
          return { ...b, elements: elems }
        }),
      )
    },
  })

  // ── Handlers ─────────────────────────────────────────────────────────────────

  async function handleCreate(event: FormEvent) {
    event.preventDefault()
    if (!createTitle.trim()) return
    setPageError(null)
    try {
      await createMutation.mutateAsync({ title: createTitle.trim() })
    } catch (caught) {
      setPageError(extractMsg(caught, 'Не удалось создать курс'))
    }
  }

  async function handleRename(event: FormEvent) {
    event.preventDefault()
    if (!renameTarget) return
    setRenameError(null)
    try {
      await renameMutation.mutateAsync({ id: renameTarget.id, title: renameTarget.title.trim() })
    } catch (caught) {
      setRenameError(extractMsg(caught, 'Не удалось переименовать курс'))
    }
  }

  async function handleAddBlock() {
    if (!activeCourseId || !newBlockTitle.trim()) return
    setPageError(null)
    try {
      await addBlockMutation.mutateAsync({ courseId: activeCourseId, title: newBlockTitle.trim() })
    } catch (caught) {
      setPageError(extractMsg(caught, 'Не удалось добавить блок'))
    }
  }

  async function handleAddElement(blockId: string) {
    const state = newElement[blockId] ?? defaultNewElement()
    if (!state.title.trim()) return
    setPageError(null)
    const fallbackName = state.fileName.trim() || state.title.trim().replace(/\s+/g, '_') || 'file'
    const payload =
      state.type === 'test'
        ? { title: state.title.trim(), type: state.type, quiz_id: state.quizId.trim(), completion_mode: state.completionMode }
        : { title: state.title.trim(), type: state.type, file_name: fallbackName, size_bytes: Number(state.sizeBytes) || 1, completion_mode: state.completionMode }
    try {
      await addElementMutation.mutateAsync({ blockId, payload })
    } catch (caught) {
      setPageError(extractMsg(caught, 'Не удалось добавить элемент'))
    }
  }

  function selectCourse(courseId: string) {
    setSelectedCourseId(courseId)
    setBlocks([])
    setNewBlockTitle('')
    setNewElement({})
    setPageError(null)
  }

  const courseBlocks = activeCourseId
    ? blocks.filter((b) => b.courseId === activeCourseId)
    : []

  const ratingsQuery = useQuery({
    queryKey: ['course-ratings', activeCourseId],
    queryFn: () => listCourseRatings(activeCourseId!, 50, 0),
    enabled: ratingsOpen && Boolean(activeCourseId),
  })

  return (
    <>
      <PageHeader
        title="Курсы"
        description="Создавайте курсы, блоки и элементы учебного контента."
      />

      <form className={styles.createForm} onSubmit={handleCreate}>
        <input
          value={createTitle}
          onChange={(e) => setCreateTitle(e.target.value)}
          placeholder="Название нового курса"
          aria-label="Название курса"
        />
        <Button type="submit" disabled={createMutation.isPending}>
          <Plus size={16} aria-hidden="true" />
          {createMutation.isPending ? 'Создание…' : 'Создать курс'}
        </Button>
      </form>

      {pageError ? <p className={styles.pageError}>{pageError}</p> : null}
      {isLoading ? <p className={styles.state}>Загрузка курсов…</p> : null}
      {error ? <p className={styles.stateError}>{error instanceof ApiError ? error.message : 'Ошибка загрузки'}</p> : null}
      {!isLoading && courses.length === 0 ? <p className={styles.state}>Курсы не найдены</p> : null}

      {courses.length > 0 ? (
        <div className={styles.builder}>
          {/* ── Sidebar ── */}
          <aside className={styles.sidebar}>
            {courses.map((course) => (
              <button
                key={course.ID}
                type="button"
                onClick={() => selectCourse(course.ID)}
                className={`${styles.courseCard} ${course.ID === activeCourse?.ID ? styles.courseCardActive : ''}`}
              >
                <span className={styles.courseCardTitle}>{course.Title}</span>
                <span className={styles.courseCardMeta}>Блоков: {course.BlocksCount}</span>
              </button>
            ))}
          </aside>

          {/* ── Workspace ── */}
          {activeCourse ? (
            <div className={styles.workspace}>
              <div className={styles.courseHead}>
                <h2 className={styles.courseTitle}>{activeCourse.Title}</h2>
                <div className={styles.courseHeadActions}>
                  <Button
                    variant="secondary"
                    onClick={() => setRatingsOpen(true)}
                  >
                    <BarChart2 size={15} aria-hidden="true" />
                    Рейтинг
                  </Button>
                  <Button
                    variant="secondary"
                    onClick={() => { setRenameTarget({ id: activeCourse.ID, title: activeCourse.Title }); setRenameError(null) }}
                  >
                    <Pencil size={15} aria-hidden="true" />
                    Переименовать
                  </Button>
                </div>
              </div>

              {/* Add block */}
              <div className={styles.addRow}>
                <input
                  className={styles.addInput}
                  value={newBlockTitle}
                  onChange={(e) => setNewBlockTitle(e.target.value)}
                  placeholder="Название блока"
                  onKeyDown={(e) => { if (e.key === 'Enter') { e.preventDefault(); void handleAddBlock() } }}
                />
                <Button onClick={() => void handleAddBlock()} disabled={addBlockMutation.isPending}>
                  <FolderPlus size={15} aria-hidden="true" />
                  Добавить блок
                </Button>
              </div>

              {courseBlocks.length === 0 ? (
                <p className={styles.state}>Блоки не добавлены в этой сессии</p>
              ) : null}

              <div className={styles.blockList}>
                {courseBlocks.map((block, blockIdx) => (
                  <BlockCard
                    key={block.id}
                    block={block}
                    blockIndex={blockIdx}
                    totalBlocks={courseBlocks.length}
                    newElementState={newElement[block.id] ?? defaultNewElement()}
                    onNewElementChange={(state) => setNewElement((prev) => ({ ...prev, [block.id]: state }))}
                    onAddElement={() => void handleAddElement(block.id)}
                    addingElement={addElementMutation.isPending}
                    onMoveUp={() => void moveBlockMutation.mutateAsync({ courseId: activeCourse.ID, from: blockIdx, to: blockIdx - 1 })}
                    onMoveDown={() => void moveBlockMutation.mutateAsync({ courseId: activeCourse.ID, from: blockIdx, to: blockIdx + 1 })}
                    onDelete={() => setDeleteBlockTarget({ courseId: activeCourse.ID, blockId: block.id, title: block.title })}
                    onMoveElement={(from, to) => void moveElementMutation.mutateAsync({ blockId: block.id, from, to })}
                    onDeleteElement={(elementId) => void deleteElementMutation.mutateAsync({ blockId: block.id, elementId })}
                    onToggleCompletion={(elementId, current) => {
                      const next: CompletionMode = current === 'manual' ? 'none' : 'manual'
                      void changeElementCompletionMode(elementId, next).then(() => {
                        setBlocks((prev) =>
                          prev.map((b) =>
                            b.id !== block.id
                              ? b
                              : { ...b, elements: b.elements.map((e) => e.id === elementId ? { ...e, completionMode: next } : e) },
                          ),
                        )
                      })
                    }}
                  />
                ))}
              </div>
            </div>
          ) : null}
        </div>
      ) : null}

      {/* ── Rename modal ── */}
      <Modal open={renameTarget !== null} onClose={() => setRenameTarget(null)} title="Переименовать курс" size="sm">
        <form className={styles.modalForm} onSubmit={handleRename} noValidate>
          <FormField label="Новое название" htmlFor="rc-title" required>
            <input
              id="rc-title"
              value={renameTarget?.title ?? ''}
              onChange={(e) => setRenameTarget((t) => t ? { ...t, title: e.target.value } : null)}
              required
              autoFocus
            />
          </FormField>
          {renameError ? <p className={styles.formError}>{renameError}</p> : null}
          <div className={styles.formActions}>
            <Button variant="secondary" type="button" onClick={() => setRenameTarget(null)} disabled={renameMutation.isPending}>Отмена</Button>
            <Button type="submit" disabled={renameMutation.isPending}>{renameMutation.isPending ? 'Сохранение…' : 'Сохранить'}</Button>
          </div>
        </form>
      </Modal>

      {/* ── Delete block confirm ── */}
      <ConfirmDialog
        open={deleteBlockTarget !== null}
        onClose={() => setDeleteBlockTarget(null)}
        onConfirm={() => { if (deleteBlockTarget) void deleteBlockMutation.mutateAsync({ courseId: deleteBlockTarget.courseId, blockId: deleteBlockTarget.blockId }) }}
        title="Удалить блок"
        message={`Блок «${deleteBlockTarget?.title ?? ''}» и все его элементы будут удалены.`}
        confirmLabel="Удалить"
        danger
        isPending={deleteBlockMutation.isPending}
      />

      {/* ── Course ratings modal ── */}
      <Modal
        open={ratingsOpen}
        onClose={() => setRatingsOpen(false)}
        title={`Рейтинг: ${activeCourse?.Title ?? ''}`}
        size="lg"
      >
        {ratingsQuery.isPending && <p className={styles.state}>Загрузка…</p>}
        {ratingsQuery.isError && <p className={styles.stateError}>Не удалось загрузить рейтинг.</p>}
        {ratingsQuery.data && ratingsQuery.data.length === 0 && (
          <p className={styles.state}>Нет данных о прогрессе студентов.</p>
        )}
        {ratingsQuery.data && ratingsQuery.data.length > 0 && (
          <div className={styles.ratingsTableWrap}>
            <table className={styles.ratingsTable}>
              <thead>
                <tr>
                  <th>Аккаунт</th>
                  <th>Зачисление</th>
                  <th>Прогресс</th>
                  <th>Выполнено</th>
                </tr>
              </thead>
              <tbody>
                {ratingsQuery.data.map((r) => (
                  <tr key={r.EnrollmentID}>
                    <td>
                      <code className={styles.ratingsUuid} title={r.AccountID}>
                        {r.AccountID.slice(0, 8)}…
                      </code>
                    </td>
                    <td>
                      <code className={styles.ratingsUuid} title={r.EnrollmentID}>
                        {r.EnrollmentID.slice(0, 8)}…
                      </code>
                    </td>
                    <td>
                      <div className={styles.ratingsProgress}>
                        <div className={styles.ratingsTrack}>
                          <div
                            className={styles.ratingsFill}
                            style={{ width: `${Math.min(100, r.CompletionPercent)}%` }}
                          />
                        </div>
                        <span>{r.CompletionPercent}%</span>
                      </div>
                    </td>
                    <td className={styles.ratingsCompleted}>
                      {r.CompletedItems} / {r.TotalItems}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Modal>
    </>
  )
}

// ── Quiz-create helpers (used inside BlockCard) ────────────────────────────────

type QuizSourceDraft = {
  _k: string
  bank_id: string
  criteria_type: 'random' | 'manual'
  question_count: number
  question_ids: string
}

function newQuizSource(): QuizSourceDraft {
  return { _k: Math.random().toString(36).slice(2), bank_id: '', criteria_type: 'random', question_count: 5, question_ids: '' }
}

function quizSourcesToPayload(sources: QuizSourceDraft[]): CreateQuizPayload['sources'] {
  return sources.map((s) => ({
    bank_id: s.bank_id,
    criteria_type: s.criteria_type,
    question_count: s.criteria_type === 'random' ? s.question_count : undefined,
    question_ids:
      s.criteria_type === 'manual'
        ? s.question_ids.split(',').map((x) => x.trim()).filter(Boolean)
        : undefined,
  }))
}

// ── BlockCard ─────────────────────────────────────────────────────────────────

function BlockCard({
  block,
  blockIndex,
  totalBlocks,
  newElementState,
  onNewElementChange,
  onAddElement,
  addingElement,
  onMoveUp,
  onMoveDown,
  onDelete,
  onMoveElement,
  onDeleteElement,
  onToggleCompletion,
}: {
  block: BlockDraft
  blockIndex: number
  totalBlocks: number
  newElementState: NewElementState
  onNewElementChange: (s: NewElementState) => void
  onAddElement: () => void
  addingElement: boolean
  onMoveUp: () => void
  onMoveDown: () => void
  onDelete: () => void
  onMoveElement: (from: number, to: number) => void
  onDeleteElement: (id: string) => void
  onToggleCompletion: (id: string, current: CompletionMode) => void
}) {
  // ── Quiz creation state ──────────────────────────────────────────────────────
  const [quizModalOpen, setQuizModalOpen] = useState(false)
  const [quizFormError, setQuizFormError] = useState<string | null>(null)
  const [quizTitle, setQuizTitle] = useState('')
  const [quizAttempts, setQuizAttempts] = useState(1)
  const [quizTimeLimit, setQuizTimeLimit] = useState(0)
  const [quizShuffle, setQuizShuffle] = useState(false)
  const [quizSources, setQuizSources] = useState<QuizSourceDraft[]>([newQuizSource()])

  const banksQuery = useQuery({
    queryKey: ['banks'],
    queryFn: () => listBanks({ limit: 200 }),
    enabled: quizModalOpen,
  })
  const banks = banksQuery.data ?? []

  const createQuizMut = useMutation({
    mutationFn: (payload: CreateQuizPayload) => createQuiz(payload),
    onSuccess: (data) => {
      onNewElementChange({ ...newElementState, quizId: data.id })
      setQuizModalOpen(false)
    },
    onError: (e) => setQuizFormError(e instanceof ApiError ? e.message : 'Не удалось создать тест'),
  })

  function openQuizModal() {
    setQuizTitle('')
    setQuizAttempts(1)
    setQuizTimeLimit(0)
    setQuizShuffle(false)
    setQuizSources([newQuizSource()])
    setQuizFormError(null)
    setQuizModalOpen(true)
  }

  function handleCreateQuiz() {
    setQuizFormError(null)
    if (quizSources.find((s) => !s.bank_id)) {
      setQuizFormError('Выберите банк для каждого источника')
      return
    }
    createQuizMut.mutate({
      title: quizTitle.trim(),
      max_attempts: quizAttempts,
      time_limit_seconds: quizTimeLimit,
      shuffle_questions: quizShuffle,
      sources: quizSourcesToPayload(quizSources),
    })
  }

  function updateSource(k: string, patch: Partial<QuizSourceDraft>) {
    setQuizSources((prev) => prev.map((s) => (s._k === k ? { ...s, ...patch } : s)))
  }

  return (
    <article className={styles.blockCard}>
      <div className={styles.blockHead}>
        <strong className={styles.blockTitle}>{block.title}</strong>
        <div className={styles.blockActions}>
          <Button variant="secondary" onClick={onMoveUp} disabled={blockIndex === 0} title="Вверх"><ArrowUp size={14} /></Button>
          <Button variant="secondary" onClick={onMoveDown} disabled={blockIndex === totalBlocks - 1} title="Вниз"><ArrowDown size={14} /></Button>
          <Button variant="secondary" onClick={onDelete} title="Удалить блок"><Trash2 size={14} /></Button>
        </div>
      </div>

      {/* Add element row */}
      <div className={styles.addElementRow}>
        <input
          className={styles.addInput}
          value={newElementState.title}
          onChange={(e) => onNewElementChange({ ...newElementState, title: e.target.value })}
          placeholder="Название элемента"
        />
        <select
          className={styles.typeSelect}
          value={newElementState.type}
          onChange={(e) => onNewElementChange({ ...newElementState, type: e.target.value as ElementType })}
        >
          <option value="lecture_material">Лекция</option>
          <option value="download_file">Файл</option>
          <option value="test">Тест</option>
        </select>
        <Button onClick={onAddElement} disabled={addingElement}>
          <Plus size={14} />
          Добавить
        </Button>
      </div>

      {newElementState.type === 'test' ? (
        <>
          <div className={styles.quizFieldRow}>
            <input
              className={styles.addInput}
              value={newElementState.quizId}
              onChange={(e) => onNewElementChange({ ...newElementState, quizId: e.target.value })}
              placeholder="ID теста (quiz_id)"
            />
            <Button variant="secondary" onClick={openQuizModal}>
              <Plus size={13} aria-hidden="true" /> Создать тест
            </Button>
          </div>
          {newElementState.quizId && (
            <p className={styles.quizIdHint}>
              ID: <code>{newElementState.quizId}</code>
            </p>
          )}
        </>
      ) : (
        <div className={styles.addElementRow}>
          <input
            className={styles.addInput}
            value={newElementState.fileName}
            onChange={(e) => onNewElementChange({ ...newElementState, fileName: e.target.value })}
            placeholder="Имя файла (напр. lecture.pdf)"
          />
          <input
            className={`${styles.addInput} ${styles.sizeInput}`}
            type="number"
            min={1}
            value={newElementState.sizeBytes}
            onChange={(e) => onNewElementChange({ ...newElementState, sizeBytes: Number(e.target.value) || 1 })}
            placeholder="Байт"
          />
        </div>
      )}

      {/* ── Quiz creation modal ── */}
      <Modal open={quizModalOpen} onClose={() => setQuizModalOpen(false)} title="Создать тест" size="md">
        <div className={styles.quizForm}>
          <FormField label="Название теста" htmlFor={`qt-${block.id}`} required>
            <input
              id={`qt-${block.id}`}
              className={styles.quizInput}
              value={quizTitle}
              onChange={(e) => setQuizTitle(e.target.value)}
              placeholder="Введите название"
            />
          </FormField>
          <div className={styles.quizRow}>
            <FormField label="Попыток (0 = без ограничений)" htmlFor={`qa-${block.id}`}>
              <input
                id={`qa-${block.id}`}
                type="number"
                min={0}
                className={styles.quizInput}
                value={quizAttempts}
                onChange={(e) => setQuizAttempts(Number(e.target.value))}
              />
            </FormField>
            <FormField label="Лимит времени, сек (0 = без)" htmlFor={`ql-${block.id}`}>
              <input
                id={`ql-${block.id}`}
                type="number"
                min={0}
                className={styles.quizInput}
                value={quizTimeLimit}
                onChange={(e) => setQuizTimeLimit(Number(e.target.value))}
              />
            </FormField>
          </div>
          <label className={styles.quizCheck}>
            <input type="checkbox" checked={quizShuffle} onChange={(e) => setQuizShuffle(e.target.checked)} />
            Перемешивать вопросы
          </label>

          <div className={styles.quizSourcesBox}>
            <div className={styles.quizSourcesHead}>
              <span className={styles.quizSourcesLabel}>Источники вопросов</span>
              <Button
                variant="secondary"
                onClick={() => setQuizSources((prev) => [...prev, newQuizSource()])}
              >
                <Plus size={12} /> Добавить
              </Button>
            </div>
            {quizSources.map((s, i) => (
              <div key={s._k} className={styles.quizSourceRow}>
                <span className={styles.quizSourceNum}>{i + 1}.</span>
                <div className={styles.quizSourceFields}>
                  <select
                    className={styles.quizSelect}
                    value={s.bank_id}
                    onChange={(e) => updateSource(s._k, { bank_id: e.target.value })}
                  >
                    <option value="">— Банк —</option>
                    {banks.map((b) => (
                      <option key={b.ID} value={b.ID}>{b.Title} ({b.QuestionsCount})</option>
                    ))}
                  </select>
                  <select
                    className={styles.quizSelect}
                    value={s.criteria_type}
                    onChange={(e) => updateSource(s._k, { criteria_type: e.target.value as 'random' | 'manual' })}
                  >
                    <option value="random">Случайные N</option>
                    <option value="manual">Конкретные ID</option>
                  </select>
                  {s.criteria_type === 'random' ? (
                    <input
                      type="number"
                      min={1}
                      className={styles.quizInput}
                      value={s.question_count}
                      onChange={(e) => updateSource(s._k, { question_count: Number(e.target.value) })}
                    />
                  ) : (
                    <input
                      className={styles.quizInput}
                      value={s.question_ids}
                      onChange={(e) => updateSource(s._k, { question_ids: e.target.value })}
                      placeholder="ID через запятую"
                    />
                  )}
                </div>
                <button
                  type="button"
                  className={styles.quizRemoveBtn}
                  onClick={() => setQuizSources((prev) => prev.filter((x) => x._k !== s._k))}
                  aria-label="Удалить источник"
                >
                  <Trash2 size={13} />
                </button>
              </div>
            ))}
          </div>

          {quizFormError && <p className={styles.quizError}>{quizFormError}</p>}
          <div className={styles.quizActions}>
            <Button variant="secondary" onClick={() => setQuizModalOpen(false)}>Отмена</Button>
            <Button
              disabled={!quizTitle.trim() || quizSources.length === 0 || createQuizMut.isPending}
              onClick={handleCreateQuiz}
            >
              {createQuizMut.isPending ? 'Создание…' : 'Создать и вставить ID'}
            </Button>
          </div>
        </div>
      </Modal>

      {block.elements.length === 0 ? (
        <p className={styles.emptyElements}>Элементы не добавлены</p>
      ) : (
        <ul className={styles.elementList}>
          {block.elements.map((el, elIdx) => (
            <ElementRow
              key={el.id}
              element={el}
              elementIndex={elIdx}
              totalElements={block.elements.length}
              onMoveUp={() => onMoveElement(elIdx, elIdx - 1)}
              onMoveDown={() => onMoveElement(elIdx, elIdx + 1)}
              onDelete={() => onDeleteElement(el.id)}
              onToggleCompletion={() => onToggleCompletion(el.id, el.completionMode)}
            />
          ))}
        </ul>
      )}
    </article>
  )
}

// ── ElementRow ────────────────────────────────────────────────────────────────

function ElementRow({
  element,
  elementIndex,
  totalElements,
  onMoveUp,
  onMoveDown,
  onDelete,
  onToggleCompletion,
}: {
  element: ElementDraft
  elementIndex: number
  totalElements: number
  onMoveUp: () => void
  onMoveDown: () => void
  onDelete: () => void
  onToggleCompletion: () => void
}) {
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [uploading, setUploading] = useState(false)
  const [uploadError, setUploadError] = useState<string | null>(null)
  const [downloading, setDownloading] = useState(false)

  const hasFile = element.type !== 'test'

  async function handleUpload(file: File) {
    setUploading(true)
    setUploadError(null)
    try {
      await uploadElementContent(element.id, file)
    } catch (caught) {
      setUploadError(caught instanceof ApiError ? caught.message : 'Ошибка загрузки')
    } finally {
      setUploading(false)
    }
  }

  async function handleDownload() {
    setDownloading(true)
    try {
      const info = await getElementDownloadURL(element.id)
      window.open(info.download_url, '_blank', 'noopener,noreferrer')
    } catch {
      // silently ignore — presigned URL fetch may fail if no file uploaded yet
    } finally {
      setDownloading(false)
    }
  }

  const TYPE_LABELS: Record<ElementType, string> = {
    lecture_material: 'Лекция',
    download_file: 'Файл',
    test: 'Тест',
  }

  return (
    <li className={styles.elementRow}>
      <div className={styles.elementInfo}>
        <span className={styles.elementTitle}>{element.title}</span>
        <span className={styles.elementMeta}>
          {TYPE_LABELS[element.type]} · {element.completionMode === 'manual' ? 'с отметкой' : 'без отметки'}
        </span>
        {uploadError ? <span className={styles.elementError}>{uploadError}</span> : null}
      </div>
      <div className={styles.elementActions}>
        {hasFile ? (
          <>
            <input
              ref={fileInputRef}
              type="file"
              className={styles.fileInput}
              aria-label="Загрузить файл"
              onChange={(e) => {
                const file = e.target.files?.[0]
                if (file) void handleUpload(file)
                e.target.value = ''
              }}
            />
            <Button
              variant="secondary"
              onClick={() => fileInputRef.current?.click()}
              disabled={uploading}
              title="Загрузить файл"
            >
              <Upload size={13} aria-hidden="true" />
              {uploading ? '…' : ''}
            </Button>
            <Button
              variant="secondary"
              onClick={() => void handleDownload()}
              disabled={downloading}
              title="Скачать файл"
            >
              <Download size={13} aria-hidden="true" />
            </Button>
          </>
        ) : null}
        <Button variant="secondary" onClick={onToggleCompletion} title="Переключить отметку завершения">
          <Pencil size={13} aria-hidden="true" />
        </Button>
        <Button variant="secondary" onClick={onMoveUp} disabled={elementIndex === 0} title="Вверх">
          <ArrowUp size={13} aria-hidden="true" />
        </Button>
        <Button variant="secondary" onClick={onMoveDown} disabled={elementIndex === totalElements - 1} title="Вниз">
          <ArrowDown size={13} aria-hidden="true" />
        </Button>
        <Button variant="secondary" onClick={onDelete} title="Удалить элемент">
          <Trash2 size={13} aria-hidden="true" />
        </Button>
      </div>
    </li>
  )
}
