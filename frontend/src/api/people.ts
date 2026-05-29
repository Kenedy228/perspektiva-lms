import type {
  CreatePersonRequest,
  PersonDetailedView,
  PersonShortView,
  ProfileRequest,
  RenamePersonRequest,
} from '../types/person'
import { request } from './client'

export function listPersonsByLastName(lastName: string, limit: number, offset: number) {
  const params = new URLSearchParams({ last_name: lastName, limit: String(limit), offset: String(offset) })
  return request<PersonShortView[]>(`/persons?${params.toString()}`)
}

export function listPersonsBySnils(snils: string, limit: number, offset: number) {
  const params = new URLSearchParams({ snils, limit: String(limit), offset: String(offset) })
  return request<PersonShortView[]>(`/persons?${params.toString()}`)
}

export function getPerson(id: string) {
  return request<PersonDetailedView>(`/persons/${id}`)
}

export function createPerson(payload: CreatePersonRequest) {
  return request<{ id: string }>('/persons', { method: 'POST', body: payload })
}

export function renamePerson(id: string, payload: RenamePersonRequest) {
  return request<{ id: string }>(`/persons/${id}`, { method: 'PATCH', body: payload })
}

export function replacePersonProfile(id: string, payload: ProfileRequest) {
  return request<{ id: string }>(`/persons/${id}/profile`, { method: 'PUT', body: payload })
}

export function detachPersonProfile(id: string) {
  return request<{ id: string }>(`/persons/${id}/profile`, { method: 'DELETE' })
}

export function assignPersonOrganization(id: string, organization_id: string) {
  return request<{ id: string }>(`/persons/${id}/organization`, {
    method: 'PUT',
    body: { organization_id },
  })
}

export function removePersonOrganization(id: string) {
  return request<{ id: string }>(`/persons/${id}/organization`, { method: 'DELETE' })
}

export function deletePerson(id: string) {
  return request<void>(`/persons/${id}`, { method: 'DELETE' })
}
