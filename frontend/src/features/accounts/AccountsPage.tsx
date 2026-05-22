import { PageHeader } from '../../components/ui/PageHeader'
import { PlaceholderList } from '../shared/PlaceholderList'

export function AccountsPage() {
  return (
    <>
      <PageHeader title="Учетные записи" description="Администраторы управляют жизненным циклом, ролями и паролями." />
      <PlaceholderList items={['Создать учетную запись', 'Заблокировать учетную запись', 'Изменить пароль']} />
    </>
  )
}
