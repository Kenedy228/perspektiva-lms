import { type ComponentPropsWithoutRef, type ElementType, type ReactNode } from 'react'
import styles from './Button.module.css'

type ButtonProps<T extends ElementType> = {
  as?: T
  children: ReactNode
  variant?: 'primary' | 'secondary'
} & Omit<ComponentPropsWithoutRef<T>, 'as' | 'children'>

export function Button<T extends ElementType = 'button'>({
  as,
  children,
  className,
  variant = 'primary',
  ...props
}: ButtonProps<T>) {
  const Component = as ?? 'button'
  const classes = [styles.button, styles[variant], className].filter(Boolean).join(' ')

  return (
    <Component className={classes} {...props}>
      {children}
    </Component>
  )
}
