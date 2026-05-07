import { lazy } from 'react';
import { Route } from 'react-router-dom';
import { AdminRoute } from './AuthRoute';
import { withGuard, withSuspense } from './routeUtils';
import AdminLayout from '../components/features/layouts/AdminLayout';

const Dashboard = lazy(() => import('../pages/admin/Dashboard'));
const ContestsManagement = lazy(() => import('../pages/admin/Contests'));
const RbacManagement = lazy(() => import('../pages/admin/Rbac'));
const ChallengesManagement = lazy(() => import('../pages/admin/Challenges'));
const FilesManagement = lazy(() => import('../pages/admin/Files.jsx'));
const SystemSettings = lazy(() => import('../pages/admin/System'));
const OAuthProvidersManagement = lazy(() => import('../pages/admin/OAuthProviders'));
const SmtpManagement = lazy(() => import('../pages/admin/Smtp'));
const CronJobsManagement = lazy(() => import('../pages/admin/CronJobs'));
const WebhookManagement = lazy(() => import('../pages/admin/Webhook'));
const BrandingManagement = lazy(() => import('../pages/admin/Branding'));
const TaskManagement = lazy(() => import('../pages/admin/Tasks'));
const AdminLogs = lazy(() => import('../pages/admin/Logs'));
const AdminImagesManagement = lazy(() => import('../pages/admin/Images.jsx'));
const AdminVictims = lazy(() => import('../pages/admin/Victims'));
const AdminGenerators = lazy(() => import('../pages/admin/Generators'));
const Settings = lazy(() => import('../pages/user/Settings'));

const ADMIN_ROUTES = [
  { path: 'dashboard', Component: Dashboard, apiRoute: 'GET /admin/system/status' },
  { path: 'contests', Component: ContestsManagement, apiRoute: 'GET /admin/contests' },
  { path: 'rbac', Component: RbacManagement, apiRoute: 'GET /admin/roles' },
  { path: 'challenges', Component: ChallengesManagement, apiRoute: 'GET /admin/challenges' },
  { path: 'oauth', Component: OAuthProvidersManagement, apiRoute: 'GET /admin/oauth' },
  { path: 'smtp', Component: SmtpManagement, apiRoute: 'GET /admin/smtp' },
  { path: 'cronjobs', Component: CronJobsManagement, apiRoute: 'GET /admin/cronjobs' },
  { path: 'webhook', Component: WebhookManagement, apiRoute: 'GET /admin/webhook' },
  { path: 'branding', Component: BrandingManagement, apiRoute: 'GET /admin/branding' },
  { path: 'files', Component: FilesManagement, apiRoute: 'GET /admin/files' },
  { path: 'tasks', Component: TaskManagement, apiRoute: 'GET /admin/tasks' },
  { path: 'system', Component: SystemSettings, apiRoute: 'GET /admin/system/config' },
  { path: 'logs', Component: AdminLogs, apiRoute: 'GET /admin/logs' },
  { path: 'victims', Component: AdminVictims, apiRoute: 'GET /admin/victims' },
  { path: 'generators', Component: AdminGenerators, apiRoute: 'GET /admin/generators' },
  { path: 'images', Component: AdminImagesManagement, apiRoute: 'GET /admin/images' },
];

export function AdminRoutes() {
  return (
    <Route path="/admin" element={withGuard(<AdminLayout />, AdminRoute)}>
      <Route path="settings" element={withSuspense(Settings)} />
      {ADMIN_ROUTES.map(({ path, Component, apiRoute }) => (
        <Route key={path} path={path} element={withGuard(withSuspense(Component), AdminRoute, { apiRoute })} />
      ))}
    </Route>
  );
}
