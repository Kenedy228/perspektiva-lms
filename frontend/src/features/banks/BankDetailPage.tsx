import { useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, Pencil, Plus, SquarePen, Trash2, Unlink } from 'lucide-react'
import {
  addQuestionsToBank,
  createQuestion,
  deleteBank,
  deleteQuestion,
  getBank,
  removeQuestionsFromBank,
  renameBank,
  updateQuestionContent,
  updateQuestionTitle,
} from '../../api/banks'
import { Badge } from '../../components/ui/Badge'
import { Button } from '../../components/ui/Button'
import { ConfirmDialog } from '../../components/ui/ConfirmDialog'
import { FormField } from '../../components/ui/FormField'
import { Modal } from '../../components/ui/Modal'
import { ApiError } from '../../types/api'
import type { BankQuestionView } from '../../types/banks'
import styles from './BankDetailPage.module.css'

// ── Local draft types ──────────────────────────────────────────────────────────

type QType = 'selectable' | 'sequence' | 'matching' | 'short'

type SelOpt = { _k: string; text: string; is_correct: boolean }
type SeqOpt = { _k: string; text: string }
type MatPair = { _k: string; prompt: string; match: string }
type ShortV = { _k: string; text: string }

type Draft = {
  type: QType
  title: string
  selOpts: SelOpt[]
  seqOpts: SeqOpt[]
  matPairs: MatPair[]
  shortVars: ShortV[]
}

function uid() {
  return Math.random().toString(36).slice(2)
}

function emptyDraft(type: QType = 'selectable'): Draft {
  return {
    type,
    title: '',
    selOpts: [
      { _k: uid(), text: '', is_correct: false },
      { _k: uid(), text: '', is_correct: false },
    ],
    seqOpts: [{ _k: uid(), text: '' }, { _k: uid(), text: '' }],
    matPairs: [{ _k: uid(), prompt: '', match: '' }],
    shortVars: [{ _k: uid(), text: '' }],
  }
}

function questionToDraft(q: BankQuestionView): Draft {
  return {
    type: q.Type,
    title: q.Title,
    selOpts:
      q.SelectableOptions?.map((o) => ({ _k: uid(), text: o.Value, is_correct: o.IsCorrect })) ?? [],
    seqOpts: q.SequenceOptions?.map((o) => ({ _k: uid(), text: o.Value })) ?? [],
    matPairs:
      q.MatchingPairs?.map((p) => ({ _k: uid(), prompt: p.PromptText, match: p.MatchText })) ?? [],
    shortVars: q.ShortVariants?.map((v) => ({ _k: uid(), text: v.Value })) ?? [],
  }
}

function draftToPayload(d: Draft) {
  return {
    type: d.type,
    title: d.title,
    selectable_options:
      d.type === 'selectable'
        ? d.selOpts.map((o) => ({ text: o.text, is_correct: o.is_correct }))
        : undefined,
    sequence_options:
      d.type === 'sequence' ? d.seqOpts.map((o) => ({ text: o.text })) : undefined,
    matching_pairs:
      d.type === 'matching'
        ? d.matPairs.map((p) => ({ prompt: p.prompt, match: p.match }))
        : undefined,
    short_variants:
      d.type === 'short' ? d.shortVars.map((v) => ({ text: v.text })) : undefined,
  }
}

// ── Type metadata ──────────────────────────────────────────────────────────────

const TYPE_LABELS: Record<QType, string> = {
  selectable: 'Выбор',
  sequence: 'Порядок',
  matching: 'Соответствие',
  short: 'Текст',
}

const TYPE_OPTIONS: { value: QType; label: string }[] = [
  { value: 'selectable', label: 'Выбор вариантов' },
  { value: 'sequence', label: 'Последовательность' },
  { value: 'matching', label: 'Соответствие' },
  { value: 'short', label: 'Открытый ответ' },
]

// ── Sub-editors for question content ──────────────────────────────────────────

function SelectableEditor({
  opts,
  onChange,
}: {
  opts: SelOpt[]
  onChange: (v: SelOpt[]) => void
}) {
  function update(k: string, patch: Partial<SelOpt>) {
    onChange(opts.map((o) => (o._k === k ? { ...o, ...patch } : o)))
  }
  function remove(k: string) {
    onChange(opts.filter((o) => o._k !== k))
  }
  return (
    <div className={styles.optList}>
      {opts.map((o) => (
        <div key={o._k} className={styles.optRow}>
          <input
            className={styles.optInput}
            value={o.text}
            placeholder="Текст варианта"
            onChange={(e) => update(o._k, { text: e.target.value })}
          />
          <label className={styles.checkLabel}>
            <input
              type="checkbox"
              checked={o.is_correct}
              onChange={(e) => update(o._k, { is_correct: e.target.checked })}
            />
            Верный
          </label>
          <button
            type="button"
            className={styles.iconBtn}
            onClick={() => remove(o._k)}
            aria-label="Удалить вариант"
          >
            <Trash2 size={14} />
          </button>
        </div>
      ))}
      <Button
        type="button"
        variant="secondary"
        onClick={() => onChange([...opts, { _k: uid(), text: '', is_correct: false }])}
      >
        <Plus size={14} /> Добавить вариант
      </Button>
    </div>
  )
}

function SequenceEditor({ opts, onChange }: { opts: SeqOpt[]; onChange: (v: SeqOpt[]) => void }) {
  function update(k: string, text: string) {
    onChange(opts.map((o) => (o._k === k ? { ...o, text } : o)))
  }
  function remove(k: string) {
    onChange(opts.filter((o) => o._k !== k))
  }
  return (
    <div className={styles.optList}>
      {opts.map((o, i) => (
        <div key={o._k} className={styles.optRow}>
          <span className={styles.seqNum}>{i + 1}.</span>
          <input
            className={styles.optInput}
            value={o.text}
            placeholder="Элемент последовательности"
            onChange={(e) => update(o._k, e.target.value)}
          />
          <button
            type="button"
            className={styles.iconBtn}
            onClick={() => remove(o._k)}
            aria-label="Удалить элемент"
          >
            <Trash2 size={14} />
          </button>
        </div>
      ))}
      <Button
        type="button"
        variant="secondary"
        onClick={() => onChange([...opts, { _k: uid(), text: '' }])}
      >
        <Plus size={14} /> Добавить элемент
      </Button>
    </div>
  )
}

function MatchingEditor({
  pairs,
  onChange,
}: {
  pairs: MatPair[]
  onChange: (v: MatPair[]) => void
}) {
  function update(k: string, patch: Partial<MatPair>) {
    onChange(pairs.map((p) => (p._k === k ? { ...p, ...patch } : p)))
  }
  function remove(k: string) {
    onChange(pairs.filter((p) => p._k !== k))
  }
  return (
    <div className={styles.optList}>
      <div className={styles.matchHeader}>
        <span>Левая часть</span>
        <span>Правая часть</span>
      </div>
      {pairs.map((p) => (
        <div key={p._k} className={styles.matchRow}>
          <input
            className={styles.optInput}
            value={p.prompt}
            placeholder="Вопрос"
            onChange={(e) => update(p._k, { prompt: e.target.value })}
          />
          <span className={styles.matchArrow}>→</span>
          <input
            className={styles.optInput}
            value={p.match}
            placeholder="Ответ"
            onChange={(e) => update(p._k, { match: e.target.value })}
          />
          <button
            type="button"
            className={styles.iconBtn}
            onClick={() => remove(p._k)}
            aria-label="Удалить пару"
          >
            <Trash2 size={14} />
          </button>
        </div>
      ))}
      <Button
        type="button"
        variant="secondary"
        onClick={() => onChange([...pairs, { _k: uid(), prompt: '', match: '' }])}
      >
        <Plus size={14} /> Добавить пару
      </Button>
    </div>
  )
}

function ShortEditor({ vars, onChange }: { vars: ShortV[]; onChange: (v: ShortV[]) => void }) {
  function update(k: string, text: string) {
    onChange(vars.map((v) => (v._k === k ? { ...v, text } : v)))
  }
  function remove(k: string) {
    onChange(vars.filter((v) => v._k !== k))
  }
  return (
    <div className={styles.optList}>
      <p className={styles.hint}>Все варианты считаются правильными ответами.</p>
      {vars.map((v) => (
        <div key={v._k} className={styles.optRow}>
          <input
            className={styles.optInput}
            value={v.text}
            placeholder="Принимаемый ответ"
            onChange={(e) => update(v._k, e.target.value)}
          />
          <button
            type="button"
            className={styles.iconBtn}
            onClick={() => remove(v._k)}
            aria-label="Удалить вариант"
          >
            <Trash2 size={14} />
          </button>
        </div>
      ))}
      <Button
        type="button"
        variant="secondary"
        onClick={() => onChange([...vars, { _k: uid(), text: '' }])}
      >
        <Plus size={14} /> Добавить вариант
      </Button>
    </div>
  )
}

// ── Question draft form (create + edit content) ────────────────────────────────

function QuestionDraftForm({
  draft,
  onChange,
  showTypeSelect,
}: {
  draft: Draft
  onChange: (d: Draft) => void
  showTypeSelect: boolean
}) {
  function set<K extends keyof Draft>(k: K, v: Draft[K]) {
    onChange({ ...draft, [k]: v })
  }
  return (
    <div className={styles.draftForm}>
      {showTypeSelect && (
        <FormField label="Тип вопроса" htmlFor="q-type" required>
          <select
            id="q-type"
            className={styles.typeSelect}
            value={draft.type}
            onChange={(e) => onChange({ ...emptyDraft(e.target.value as QType), title: draft.title })}
          >
            {TYPE_OPTIONS.map((o) => (
              <option key={o.value} value={o.value}>
                {o.label}
              </option>
            ))}
          </select>
        </FormField>
      )}

      <FormField label="Формулировка вопроса" htmlFor="q-title" required>
        <input
          id="q-title"
          className={styles.titleInput}
          value={draft.title}
          onChange={(e) => set('title', e.target.value)}
          placeholder="Введите текст вопроса"
        />
      </FormField>

      {draft.type === 'selectable' && (
        <FormField label="Варианты ответа">
          <SelectableEditor opts={draft.selOpts} onChange={(v) => set('selOpts', v)} />
        </FormField>
      )}
      {draft.type === 'sequence' && (
        <FormField label="Элементы в правильном порядке">
          <SequenceEditor opts={draft.seqOpts} onChange={(v) => set('seqOpts', v)} />
        </FormField>
      )}
      {draft.type === 'matching' && (
        <FormField label="Пары соответствий">
          <MatchingEditor pairs={draft.matPairs} onChange={(v) => set('matPairs', v)} />
        </FormField>
      )}
      {draft.type === 'short' && (
        <FormField label="Принимаемые варианты ответа">
          <ShortEditor vars={draft.shortVars} onChange={(v) => set('shortVars', v)} />
        </FormField>
      )}
    </div>
  )
}

// ── Question preview (collapsed view in list) ──────────────────────────────────

function QuestionPreview({ q }: { q: BankQuestionView }) {
  if (q.Type === 'selectable' && q.SelectableOptions?.length) {
    const shown = q.SelectableOptions.slice(0, 3)
    return (
      <ul className={styles.preview}>
        {shown.map((o) => (
          <li key={o.ID} className={o.IsCorrect ? styles.correct : undefined}>
            {o.IsCorrect ? '✓' : '○'} {o.Value || '—'}
          </li>
        ))}
        {q.SelectableOptions.length > 3 && (
          <li className={styles.more}>ещё {q.SelectableOptions.length - 3}…</li>
        )}
      </ul>
    )
  }
  if (q.Type === 'sequence' && q.SequenceOptions?.length) {
    const shown = q.SequenceOptions.slice(0, 3)
    return (
      <ul className={styles.preview}>
        {shown.map((o, i) => (
          <li key={i}>{i + 1}. {o.Value || '—'}</li>
        ))}
        {q.SequenceOptions.length > 3 && (
          <li className={styles.more}>ещё {q.SequenceOptions.length - 3}…</li>
        )}
      </ul>
    )
  }
  if (q.Type === 'matching' && q.MatchingPairs?.length) {
    const shown = q.MatchingPairs.slice(0, 2)
    return (
      <ul className={styles.preview}>
        {shown.map((p) => (
          <li key={p.PromptID}>{p.PromptText || '—'} → {p.MatchText || '—'}</li>
        ))}
        {q.MatchingPairs.length > 2 && (
          <li className={styles.more}>ещё {q.MatchingPairs.length - 2}…</li>
        )}
      </ul>
    )
  }
  if (q.Type === 'short' && q.ShortVariants?.length) {
    return (
      <p className={styles.previewShort}>
        Ответы: {q.ShortVariants.map((v) => v.Value).filter(Boolean).join(', ')}
      </p>
    )
  }
  return null
}

// ── Main page ──────────────────────────────────────────────────────────────────

export function BankDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const bankQuery = useQuery({
    queryKey: ['bank', id],
    queryFn: () => getBank(id!),
    enabled: Boolean(id),
  })
  const bank = bankQuery.data

  // Modal open state
  const [renameOpen, setRenameOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [createOpen, setCreateOpen] = useState(false)
  const [editTitleQ, setEditTitleQ] = useState<BankQuestionView | null>(null)
  const [editContentQ, setEditContentQ] = useState<BankQuestionView | null>(null)
  const [removeQ, setRemoveQ] = useState<BankQuestionView | null>(null)
  const [deleteQ, setDeleteQ] = useState<BankQuestionView | null>(null)

  // Form state
  const [newBankTitle, setNewBankTitle] = useState('')
  const [newQTitle, setNewQTitle] = useState('')
  const [createDraft, setCreateDraft] = useState<Draft>(emptyDraft)
  const [editDraft, setEditDraft] = useState<Draft>(emptyDraft)
  const [formError, setFormError] = useState<string | null>(null)

  function inv() {
    return queryClient.invalidateQueries({ queryKey: ['bank', id] })
  }

  // ── Mutations ────────────────────────────────────────────────────────────────

  const renameBankMut = useMutation({
    mutationFn: (title: string) => renameBank(id!, title),
    onSuccess: async () => { await inv(); setRenameOpen(false) },
    onError: (e) => setFormError(e instanceof ApiError ? e.message : 'Не удалось переименовать'),
  })

  const deleteBankMut = useMutation({
    mutationFn: () => deleteBank(id!),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['banks'] })
      navigate('/banks')
    },
  })

  const createQMut = useMutation({
    mutationFn: async (draft: Draft) => {
      const q = await createQuestion(draftToPayload(draft))
      await addQuestionsToBank(id!, [q.id])
    },
    onSuccess: async () => { await inv(); setCreateOpen(false) },
    onError: (e) => setFormError(e instanceof ApiError ? e.message : 'Не удалось создать вопрос'),
  })

  const editTitleMut = useMutation({
    mutationFn: ({ qid, title }: { qid: string; title: string }) => updateQuestionTitle(qid, title),
    onSuccess: async () => { await inv(); setEditTitleQ(null) },
    onError: (e) => setFormError(e instanceof ApiError ? e.message : 'Не удалось сохранить'),
  })

  const editContentMut = useMutation({
    mutationFn: ({ qid, draft }: { qid: string; draft: Draft }) =>
      updateQuestionContent(qid, draftToPayload(draft)),
    onSuccess: async () => { await inv(); setEditContentQ(null) },
    onError: (e) => setFormError(e instanceof ApiError ? e.message : 'Не удалось сохранить'),
  })

  const removeMut = useMutation({
    mutationFn: (qid: string) => removeQuestionsFromBank(id!, [qid]),
    onSuccess: async () => { await inv(); setRemoveQ(null) },
  })

  const deleteQMut = useMutation({
    mutationFn: (qid: string) => deleteQuestion(qid),
    onSuccess: async () => { await inv(); setDeleteQ(null) },
  })

  // ── Render ────────────────────────────────────────────────────────────────────

  if (bankQuery.isPending) {
    return <p className={styles.state}>Загрузка…</p>
  }
  if (bankQuery.isError || !bank) {
    return <p className={styles.stateError}>Не удалось загрузить банк вопросов.</p>
  }

  return (
    <>
      {/* Navigation */}
      <Link to="/banks" className={styles.backLink}>
        <ArrowLeft size={15} aria-hidden="true" />
        Банки вопросов
      </Link>

      {/* Header */}
      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <h1 className={styles.title}>{bank.Title}</h1>
          <span className={styles.count}>{bank.Questions?.length ?? 0} вопросов</span>
        </div>
        <div className={styles.headerActions}>
          <Button
            variant="secondary"
            onClick={() => {
              setNewBankTitle(bank.Title)
              setFormError(null)
              setRenameOpen(true)
            }}
          >
            <Pencil size={15} /> Переименовать
          </Button>
          <Button variant="secondary" onClick={() => setDeleteOpen(true)}>
            <Trash2 size={15} /> Удалить банк
          </Button>
        </div>
      </div>

      {/* Add question */}
      <div className={styles.addRow}>
        <Button
          onClick={() => {
            setCreateDraft(emptyDraft())
            setFormError(null)
            setCreateOpen(true)
          }}
        >
          <Plus size={15} /> Создать вопрос
        </Button>
      </div>

      {/* Question list */}
      {!bank.Questions?.length ? (
        <p className={styles.empty}>Вопросов нет. Создайте первый.</p>
      ) : (
        <ul className={styles.questionList}>
          {bank.Questions.map((q) => (
            <li key={q.ID} className={styles.questionCard}>
              <div className={styles.questionHead}>
                <div className={styles.questionMeta}>
                  <Badge variant="neutral">{TYPE_LABELS[q.Type] ?? q.Type}</Badge>
                  <span className={styles.questionTitle}>{q.Title}</span>
                </div>
                <div className={styles.questionActions}>
                  <button
                    type="button"
                    className={styles.iconBtn}
                    title="Редактировать название"
                    onClick={() => {
                      setNewQTitle(q.Title)
                      setFormError(null)
                      setEditTitleQ(q)
                    }}
                  >
                    <Pencil size={15} />
                  </button>
                  <button
                    type="button"
                    className={styles.iconBtn}
                    title="Редактировать содержимое"
                    onClick={() => {
                      setEditDraft(questionToDraft(q))
                      setFormError(null)
                      setEditContentQ(q)
                    }}
                  >
                    <SquarePen size={15} />
                  </button>
                  <button
                    type="button"
                    className={styles.iconBtn}
                    title="Убрать из банка"
                    onClick={() => setRemoveQ(q)}
                  >
                    <Unlink size={15} />
                  </button>
                  <button
                    type="button"
                    className={`${styles.iconBtn} ${styles.danger}`}
                    title="Удалить вопрос"
                    onClick={() => setDeleteQ(q)}
                  >
                    <Trash2 size={15} />
                  </button>
                </div>
              </div>
              <QuestionPreview q={q} />
            </li>
          ))}
        </ul>
      )}

      {/* ── Rename bank modal ── */}
      <Modal open={renameOpen} onClose={() => setRenameOpen(false)} title="Переименовать банк" size="sm">
        <div className={styles.modalForm}>
          <FormField label="Название" htmlFor="rename-bank" required>
            <input
              id="rename-bank"
              className={styles.input}
              value={newBankTitle}
              onChange={(e) => setNewBankTitle(e.target.value)}
            />
          </FormField>
          {formError && <p className={styles.formError}>{formError}</p>}
          <div className={styles.formActions}>
            <Button variant="secondary" onClick={() => setRenameOpen(false)}>
              Отмена
            </Button>
            <Button
              disabled={!newBankTitle.trim() || renameBankMut.isPending}
              onClick={() => renameBankMut.mutate(newBankTitle.trim())}
            >
              {renameBankMut.isPending ? 'Сохранение…' : 'Сохранить'}
            </Button>
          </div>
        </div>
      </Modal>

      {/* ── Delete bank confirm ── */}
      <ConfirmDialog
        open={deleteOpen}
        onClose={() => setDeleteOpen(false)}
        onConfirm={() => deleteBankMut.mutate()}
        title="Удалить банк вопросов"
        message={`Банк «${bank.Title}» и все его вопросы будут удалены. Это действие необратимо.`}
        confirmLabel="Удалить"
        danger
        isPending={deleteBankMut.isPending}
      />

      {/* ── Create question modal ── */}
      <Modal open={createOpen} onClose={() => setCreateOpen(false)} title="Создать вопрос">
        <div className={styles.modalForm}>
          <QuestionDraftForm
            draft={createDraft}
            onChange={setCreateDraft}
            showTypeSelect
          />
          {formError && <p className={styles.formError}>{formError}</p>}
          <div className={styles.formActions}>
            <Button variant="secondary" onClick={() => setCreateOpen(false)}>
              Отмена
            </Button>
            <Button
              disabled={!createDraft.title.trim() || createQMut.isPending}
              onClick={() => {
                setFormError(null)
                createQMut.mutate(createDraft)
              }}
            >
              {createQMut.isPending ? 'Создание…' : 'Создать'}
            </Button>
          </div>
        </div>
      </Modal>

      {/* ── Edit question title modal ── */}
      <Modal
        open={Boolean(editTitleQ)}
        onClose={() => setEditTitleQ(null)}
        title="Редактировать название"
        size="sm"
      >
        <div className={styles.modalForm}>
          <FormField label="Название вопроса" htmlFor="edit-q-title" required>
            <input
              id="edit-q-title"
              className={styles.input}
              value={newQTitle}
              onChange={(e) => setNewQTitle(e.target.value)}
            />
          </FormField>
          {formError && <p className={styles.formError}>{formError}</p>}
          <div className={styles.formActions}>
            <Button variant="secondary" onClick={() => setEditTitleQ(null)}>
              Отмена
            </Button>
            <Button
              disabled={!newQTitle.trim() || editTitleMut.isPending}
              onClick={() => {
                setFormError(null)
                editTitleMut.mutate({ qid: editTitleQ!.ID, title: newQTitle.trim() })
              }}
            >
              {editTitleMut.isPending ? 'Сохранение…' : 'Сохранить'}
            </Button>
          </div>
        </div>
      </Modal>

      {/* ── Edit question content modal ── */}
      <Modal
        open={Boolean(editContentQ)}
        onClose={() => setEditContentQ(null)}
        title="Редактировать содержимое"
      >
        <div className={styles.modalForm}>
          <QuestionDraftForm
            draft={editDraft}
            onChange={setEditDraft}
            showTypeSelect={false}
          />
          {formError && <p className={styles.formError}>{formError}</p>}
          <div className={styles.formActions}>
            <Button variant="secondary" onClick={() => setEditContentQ(null)}>
              Отмена
            </Button>
            <Button
              disabled={!editDraft.title.trim() || editContentMut.isPending}
              onClick={() => {
                setFormError(null)
                editContentMut.mutate({ qid: editContentQ!.ID, draft: editDraft })
              }}
            >
              {editContentMut.isPending ? 'Сохранение…' : 'Сохранить'}
            </Button>
          </div>
        </div>
      </Modal>

      {/* ── Remove from bank confirm ── */}
      <ConfirmDialog
        open={Boolean(removeQ)}
        onClose={() => setRemoveQ(null)}
        onConfirm={() => removeMut.mutate(removeQ!.ID)}
        title="Убрать вопрос из банка"
        message={`Вопрос «${removeQ?.Title}» будет убран из банка, но не удалён.`}
        confirmLabel="Убрать"
        isPending={removeMut.isPending}
      />

      {/* ── Delete question confirm ── */}
      <ConfirmDialog
        open={Boolean(deleteQ)}
        onClose={() => setDeleteQ(null)}
        onConfirm={() => deleteQMut.mutate(deleteQ!.ID)}
        title="Удалить вопрос"
        message={`Вопрос «${deleteQ?.Title}» будет удалён безвозвратно.`}
        confirmLabel="Удалить"
        danger
        isPending={deleteQMut.isPending}
      />
    </>
  )
}
