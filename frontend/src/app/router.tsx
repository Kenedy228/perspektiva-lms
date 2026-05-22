import { createBrowserRouter, Navigate } from 'react-router-dom'
import { AppLayout } from '../components/layout/AppLayout'
import { LoginPage } from '../features/auth/LoginPage'
import { DashboardPage } from '../features/dashboard/DashboardPage'
import { CoursesPage } from '../features/courses/CoursesPage'
import { OrganizationsPage } from '../features/organizations/OrganizationsPage'
import { PeoplePage } from '../features/people/PeoplePage'
import { AccountsPage } from '../features/accounts/AccountsPage'
import { RouteError } from '../components/layout/RouteError'

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
    errorElement: <RouteError />,
  },
  {
    path: '/',
    element: <AppLayout />,
    errorElement: <RouteError />,
    children: [
      { index: true, element: <DashboardPage /> },
      { path: 'courses', element: <CoursesPage /> },
      { path: 'organizations', element: <OrganizationsPage /> },
      { path: 'people', element: <PeoplePage /> },
      { path: 'accounts', element: <AccountsPage /> },
    ],
  },
  {
    path: '*',
    element: <Navigate to="/" replace />,
  },
])
