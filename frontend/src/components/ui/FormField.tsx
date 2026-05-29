import { type ReactNode } from 'react'
import styles from './FormField.module.css'

type FormFieldProps = {
  label: string
  htmlFor?: string
  error?: string | null
  required?: boolean
  children: ReactNode
}

export function FormField({ label, htmlFor, error, required, children }: FormFieldProps) {
  return (
    <div className={styles.field}>
      <label htmlFor={htmlFor} className={styles.label}>
        {label}
        {required ? <span className={styles.required} aria-hidden="true"> *</span> : null}
      </label>
      {children}
      {error ? <span className={styles.error} role="alert">{error}</span> : null}
    </div>
  )
}
