import { Button } from './Button'
import { Modal } from './Modal'
import styles from './ConfirmDialog.module.css'

type ConfirmDialogProps = {
  open: boolean
  onClose: () => void
  onConfirm: () => void
  title: string
  message: string
  confirmLabel?: string
  cancelLabel?: string
  danger?: boolean
  isPending?: boolean
}

export function ConfirmDialog({
  open,
  onClose,
  onConfirm,
  title,
  message,
  confirmLabel = 'Подтвердить',
  cancelLabel = 'Отмена',
  danger = false,
  isPending = false,
}: ConfirmDialogProps) {
  return (
    <Modal open={open} onClose={onClose} title={title} size="sm">
      <p className={styles.message}>{message}</p>
      <div className={styles.actions}>
        <Button variant="secondary" onClick={onClose} disabled={isPending}>
          {cancelLabel}
        </Button>
        <button
          type="button"
          className={`${styles.confirmBtn} ${danger ? styles.danger : styles.primary}`}
          onClick={onConfirm}
          disabled={isPending}
        >
          {isPending ? 'Выполняется…' : confirmLabel}
        </button>
      </div>
    </Modal>
  )
}
