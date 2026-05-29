import styles from './Badge.module.css'

type BadgeVariant = 'success' | 'danger' | 'warning' | 'neutral'

type BadgeProps = {
  variant: BadgeVariant
  children: string
}

export function Badge({ variant, children }: BadgeProps) {
  return <span className={`${styles.badge} ${styles[variant]}`}>{children}</span>
}
