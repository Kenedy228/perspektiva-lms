import { type OrganizationRequest, type OrganizationView } from '../types/organizations'
import { request } from './client'

export function getOrganizations() {
  return request<OrganizationView[]>('/organizations')
}

export function createOrganization(payload: OrganizationRequest) {
  return request<{ id: string }>('/organizations', {
    method: 'POST',
    body: payload,
  })
}
