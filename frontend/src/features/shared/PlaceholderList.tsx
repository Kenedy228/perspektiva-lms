import styles from './PlaceholderList.module.css'

type PlaceholderListProps = {
  items: string[]
}

export function PlaceholderList({ items }: PlaceholderListProps) {
  return (
    <div className={styles.list}>
      {items.map((item) => (
        <div className={styles.item} key={item}>
          {item}
        </div>
      ))}
    </div>
  )
}
