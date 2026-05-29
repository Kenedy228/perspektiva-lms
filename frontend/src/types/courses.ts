export type CourseShortView = {
  ID: string
  Title: string
  Published: boolean
  BlocksCount: number
}

export type CourseDetailedView = {
  ID: string
  Title: string
  Blocks: CourseBlockView[]
}

export type CourseBlockView = {
  ID: string
  Title: string
  Elements: CourseElementView[]
}

export type CourseElementView = {
  ID: string
  Title: string
  Type: string
  CompletionMode: string
  QuizID: string
}

export type StudentRatingView = {
  AccountID: string
  EnrollmentID: string
  CourseID: string
  CompletionPercent: number
  CompletedItems: number
  TotalItems: number
}

export type CreateCourseRequest = {
  title: string
}

export type ElementType = 'lecture_material' | 'download_file' | 'test'
export type CompletionMode = 'none' | 'manual'

export type ElementDownloadInfo = {
  element_id: string
  download_url: string
  file_name: string
  content_type: string
}
