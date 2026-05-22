import { PageHeader } from '../../components/ui/PageHeader'
import { PlaceholderList } from '../shared/PlaceholderList'

export function PeoplePage() {
  return (
    <>
      <PageHeader title="Сотрудники" description="Профили связывают учетные записи с организациями и статистикой студентов." />
      <PlaceholderList items={['Создать сотрудника', 'Добавить профиль', 'Назначить организацию']} />
    </>
  )
}
