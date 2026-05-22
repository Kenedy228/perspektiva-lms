import { type PersonRequest, type PersonView } from '../types/person'
import { request } from './client'

export function getPeople() {
  return request<PersonView[]>('/persons')
}

export function createPerson(payload: PersonRequest) {
  return request<{ id: string }>('/persons', {
    method: 'POST',
    body: payload,
  })
}
