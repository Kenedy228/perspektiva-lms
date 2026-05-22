import { PageHeader } from '../../components/ui/PageHeader'
import styles from './DashboardPage.module.css'

const metrics = [
  { label: 'Курсы', value: '0', detail: 'Данные появятся после загрузки из backend' },
  { label: 'Студенты', value: '0', detail: 'Endpoint статистики готов' },
  { label: 'Активные учетные записи', value: '0', detail: 'Управляются администраторами' },
]

export function DashboardPage() {
  return (
    <>
      <PageHeader title="Обзор" description="Оперативная сводка по курсам, зачислениям и прогрессу студентов." />
      <section className={styles.grid}>
        {metrics.map((metric) => (
          <article className={styles.metric} key={metric.label}>
            <span>{metric.label}</span>
            <strong>{metric.value}</strong>
            <p>{metric.detail}</p>
          </article>
        ))}
      </section>
    </>
  )
}
