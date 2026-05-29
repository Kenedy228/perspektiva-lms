export type InnType = 'individual entrepreneur' | 'natural person' | 'legal entity'

export const INN_TYPE_LABELS: Record<InnType, string> = {
  'individual entrepreneur': 'Индивидуальный предприниматель',
  'natural person': 'Физическое лицо',
  'legal entity': 'Юридическое лицо',
}

export const INN_TYPES: InnType[] = [
  'individual entrepreneur',
  'natural person',
  'legal entity',
]

export type OrganizationShortView = {
  ID: string
  OrganizationName: string
}

export type OrganizationDetailedView = {
  ID: string
  OrganizationName: string
  INN: string
  INNTitle: string
}

export type CreateOrganizationRequest = {
  name: string
  inn?: string
  inn_type?: InnType
}
