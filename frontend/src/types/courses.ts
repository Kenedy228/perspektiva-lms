export type CourseShortView = {
  ID: string
  Title: string
  Published: boolean
  BlocksCount: number
}

export type CourseVersionView = {
  ID: string
  Title: string
  Status: string
}

export type CourseDetailedView = {
  ID: string
  Title: string
  Versions: CourseVersionView[]
}

export type StudentRatingView = {
  AccountID: string
  EnrollmentID: string
  CourseID: string
  VersionID: string
  CompletionPercent: number
  CompletedItems: number
  TotalItems: number
}

export type CreateCourseRequest = {
  title: string
}
