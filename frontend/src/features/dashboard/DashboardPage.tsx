import { BarChart2, BookMarked, BookOpen, Building2, ClipboardList, Database, UserCog, UsersRound } from 'lucide-react'
import { Link } from 'react-router-dom'
import { useSession } from '../auth/useSession'
import { PageHeader } from '../../components/ui/PageHeader'
import styles from './DashboardPage.module.css'

type NavTile = {
  to: string
  label: string
  description: string
  icon: React.ElementType
  adminOnly?: boolean
}

const tiles: NavTile[] = [
  { to: '/courses', label: 'Курсы', description: 'Управление курсами, блоками и материалами', icon: BookOpen },
  { to: '/banks', label: 'Банки вопросов', description: 'Создание и редактирование вопросов', icon: Database },
  { to: '/quizzes', label: 'Тесты', description: 'Настройка тестов и источников вопросов', icon: ClipboardList },
  { to: '/organizations', label: 'Организации', description: 'Список и управление организациями', icon: Building2, adminOnly: true },
  { to: '/people', label: 'Сотрудники', description: 'Профили сотрудников и их организации', icon: UsersRound, adminOnly: true },
  { to: '/accounts', label: 'Учётные записи', description: 'Управление аккаунтами и ролями', icon: UserCog, adminOnly: true },
  { to: '/enrollments', label: 'Зачисления', description: 'Зачисление студентов на курсы', icon: BookMarked, adminOnly: true },
  { to: '/statistics', label: 'Статистика', description: 'Прогресс студентов по зачислениям', icon: BarChart2, adminOnly: true },
]

export function DashboardPage() {
  const { isAdmin } = useSession()
  const visibleTiles = tiles.filter((t) => !t.adminOnly || isAdmin)

  return (
    <>
      <PageHeader
        title="Главная"
        description="Добро пожаловать в панель управления LMS."
      />
      <section className={styles.grid}>
        {visibleTiles.map((tile) => (
          <Link key={tile.to} to={tile.to} className={styles.tile}>
            <tile.icon size={24} className={styles.tileIcon} aria-hidden="true" />
            <strong className={styles.tileLabel}>{tile.label}</strong>
            <p className={styles.tileDesc}>{tile.description}</p>
          </Link>
        ))}
      </section>
    </>
  )
}
