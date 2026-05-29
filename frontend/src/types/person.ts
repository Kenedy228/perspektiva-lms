export type PersonShortView = {
  ID: string
  FullName: string
  OrganizationName: string
}

export type PersonDetailedView = {
  ID: string
  FirstName: string
  LastName: string
  MiddleName: string
  Snils: string
  JobTitle: string
  Education: string
  DateOfBirth: string
  OrganizationName: string
}

export type CreatePersonRequest = {
  first_name: string
  last_name: string
  middle_name?: string
}

export type RenamePersonRequest = {
  first_name: string
  last_name: string
  middle_name?: string
}

export type ProfileRequest = {
  snils: string
  date_of_birth: string
  job_title?: string
  education?: string
  organization_id?: string
}
