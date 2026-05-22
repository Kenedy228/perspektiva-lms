import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { PlaceholderList } from './PlaceholderList'

describe('PlaceholderList', () => {
  it('renders provided items', () => {
    render(<PlaceholderList items={['Создать учетную запись', 'Заблокировать учетную запись']} />)

    expect(screen.getByText('Создать учетную запись')).toBeInTheDocument()
    expect(screen.getByText('Заблокировать учетную запись')).toBeInTheDocument()
  })
})
