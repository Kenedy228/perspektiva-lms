export type PersonView = {
  ID: string
  FirstName: string
  LastName: string
  MiddleName?: string
  OrganizationID?: string
}

export type PersonRequest = {
  first_name: string
  last_name: string
  middle_name?: string
}
