import { createBrowserRouter, Navigate } from 'react-router-dom'
import { AppLayout } from '../components/layout/AppLayout'
import { StudentLayout } from '../components/layout/StudentLayout'
import { OrgLayout } from '../components/layout/OrgLayout'
import { AuthGuard } from '../features/auth/AuthGuard'
import { AdminOnlyGuard } from '../features/auth/AdminOnlyGuard'
import { StudentGuard } from '../features/auth/StudentGuard'
import { OrgGuard } from '../features/auth/OrgGuard'
import { LoginPage } from '../features/auth/LoginPage'
import { MyCoursesPage } from '../features/student/MyCoursesPage'
import { MyEnrollmentsPage } from '../features/student/MyEnrollmentsPage'
import { AttemptPage } from '../features/student/AttemptPage'
import { MyCourseDetailPage } from '../features/student/MyCourseDetailPage'
import { CoursesPage } from '../features/courses/CoursesPage'
import { BanksPage } from '../features/banks/BanksPage'
import { BankDetailPage } from '../features/banks/BankDetailPage'
import { CreateBankPage } from '../features/banks/CreateBankPage'
import { QuizzesPage } from '../features/quizzes/QuizzesPage'
import { QuizDetailPage } from '../features/quizzes/QuizDetailPage'
import { OrganizationsPage } from '../features/organizations/OrganizationsPage'
import { PeoplePage } from '../features/people/PeoplePage'
import { PersonDetailPage } from '../features/people/PersonDetailPage'
import { AccountsPage } from '../features/accounts/AccountsPage'
import { EnrollmentsPage } from '../features/enrollments/EnrollmentsPage'
import { StatisticsPage } from '../features/statistics/StatisticsPage'
import { DashboardPage } from '../features/dashboard/DashboardPage'
import { EnrollmentAttemptsPage } from '../features/attempts/EnrollmentAttemptsPage'
import { AttemptAdminDetailPage } from '../features/attempts/AttemptAdminDetailPage'
import { OrgStatisticsPage } from '../features/org/OrgStatisticsPage'
import { RouteError } from '../components/layout/RouteError'

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
    errorElement: <RouteError />,
  },
  {
    path: '/',
    element: <AuthGuard />,
    errorElement: <RouteError />,
    children: [
      {
        path: '/',
        element: <AppLayout />,
        errorElement: <RouteError />,
        children: [
          { index: true, element: <DashboardPage /> },
          { path: 'courses', element: <CoursesPage /> },
          { path: 'banks', element: <BanksPage /> },
          { path: 'banks/new', element: <CreateBankPage /> },
          { path: 'banks/:id', element: <BankDetailPage /> },
          { path: 'quizzes', element: <QuizzesPage /> },
          { path: 'quizzes/:id', element: <QuizDetailPage /> },
          { path: 'attempts/:id', element: <AttemptAdminDetailPage /> },
          {
            path: '/',
            element: <AdminOnlyGuard />,
            children: [
              { path: 'organizations', element: <OrganizationsPage /> },
              { path: 'people', element: <PeoplePage /> },
              { path: 'people/:id', element: <PersonDetailPage /> },
              { path: 'accounts', element: <AccountsPage /> },
              { path: 'enrollments', element: <EnrollmentsPage /> },
              { path: 'enrollments/:id/attempts', element: <EnrollmentAttemptsPage /> },
              { path: 'statistics', element: <StatisticsPage /> },
            ],
          },
        ],
      },
    ],
  },
  {
    path: '/my',
    element: <StudentGuard />,
    errorElement: <RouteError />,
    children: [
      {
        path: '/my',
        element: <StudentLayout />,
        errorElement: <RouteError />,
        children: [
          { index: true, element: <Navigate to="/my/courses" replace /> },
          { path: 'courses', element: <MyCoursesPage /> },
          { path: 'courses/:id', element: <MyCourseDetailPage /> },
          { path: 'enrollments', element: <MyEnrollmentsPage /> },
          { path: 'attempts/:id', element: <AttemptPage /> },
        ],
      },
    ],
  },
  {
    path: '/org',
    element: <OrgGuard />,
    errorElement: <RouteError />,
    children: [
      {
        path: '/org',
        element: <OrgLayout />,
        errorElement: <RouteError />,
        children: [
          { index: true, element: <Navigate to="/org/statistics" replace /> },
          { path: 'statistics', element: <OrgStatisticsPage /> },
        ],
      },
    ],
  },
  {
    path: '*',
    element: <Navigate to="/" replace />,
  },
])
