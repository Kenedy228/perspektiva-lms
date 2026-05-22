import { PageHeader } from '../../components/ui/PageHeader'
import { PlaceholderList } from '../shared/PlaceholderList'

export function OrganizationsPage() {
  return (
    <>
      <PageHeader title="Организации" description="Записи организаций и ИНН управляются администраторами." />
      <PlaceholderList items={['Создать организацию', 'Посмотреть связанных сотрудников', 'Обновить ИНН организации']} />
    </>
  )
}
