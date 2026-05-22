export type OrganizationView = {
  ID: string
  Name: string
  INN: string
  INNType: string
}

export type OrganizationRequest = {
  name: string
  inn?: string
  inn_type?: string
}
