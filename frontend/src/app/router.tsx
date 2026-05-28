import { createBrowserRouter, Navigate } from 'react-router-dom'
import { AppLayout } from '../components/layout/AppLayout'
import { LoginPage } from '../features/auth/LoginPage'
import { CoursesPage } from '../features/courses/CoursesPage'
import { BanksPage } from '../features/banks/BanksPage'
import { CreateBankPage } from '../features/banks/CreateBankPage'
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
      { index: true, element: <Navigate to="/courses" replace /> },
      { path: 'courses', element: <CoursesPage /> },
      { path: 'banks', element: <BanksPage /> },
      { path: 'banks/new', element: <CreateBankPage /> },
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
