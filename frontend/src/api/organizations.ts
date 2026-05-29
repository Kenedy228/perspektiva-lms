import type {
  CreateOrganizationRequest,
  InnType,
  OrganizationDetailedView,
  OrganizationShortView,
} from '../types/organizations'
import { request } from './client'

export function listOrganizationsByName(search: string, limit: number, offset: number) {
  const params = new URLSearchParams({ name: search, limit: String(limit), offset: String(offset) })
  return request<OrganizationShortView[]>(`/organizations?${params.toString()}`)
}

export function listOrganizationsByINN(search: string, limit: number, offset: number) {
  const params = new URLSearchParams({ inn: search, limit: String(limit), offset: String(offset) })
  return request<OrganizationShortView[]>(`/organizations?${params.toString()}`)
}

export function getOrganization(id: string) {
  return request<OrganizationDetailedView>(`/organizations/${id}`)
}

export function createOrganization(payload: CreateOrganizationRequest) {
  return request<{ id: string }>('/organizations', { method: 'POST', body: payload })
}

export function renameOrganization(id: string, name: string) {
  return request<{ id: string }>(`/organizations/${id}`, { method: 'PATCH', body: { name } })
}

export function changeOrganizationINN(id: string, inn: string, inn_type: InnType) {
  return request<{ id: string }>(`/organizations/${id}/inn`, {
    method: 'PATCH',
    body: { inn, inn_type },
  })
}

export function deleteOrganization(id: string) {
  return request<void>(`/organizations/${id}`, { method: 'DELETE' })
}
