import { lazy } from 'react';
import { Route } from 'react-router-dom';
import { AdminRoute } from './AuthRoute';
import { withGuard, withSuspense } from './routeUtils';
import AdminContestsLayout from '../components/features/layouts/AdminContestsLayout';

const AdminContestDetail = lazy(() => import('../pages/admin/contests/index.jsx'));
const AdminContestChallenges = lazy(() => import('../pages/admin/contests/challenges'));
const AdminContestScoreboard = lazy(() => import('../pages/admin/contests/scoreboard'));
const AdminContestTeams = lazy(() => import('../pages/admin/contests/teams'));
const AdminContestSettings = lazy(() => import('../pages/admin/contests/settings'));
const TeamDetails = lazy(() => import('../pages/admin/contests/team-details'));
const AdminContestNotices = lazy(() => import('../pages/admin/contests/notice'));
const AdminContestImagesPull = lazy(() => import('../pages/admin/contests/images-pull.jsx'));
const ContestContainers = lazy(() => import('../pages/admin/contests/containers'));
const AdminContestCheats = lazy(() => import('../pages/admin/contests/cheats'));
const AdminContestGenerators = lazy(() => import('../pages/admin/contests/generators.jsx'));

const ADMIN_CONTEST_ROUTES = [
  { path: 'challenges', Component: AdminContestChallenges, apiRoute: 'GET /admin/contests/:contestID/challenges' },
  { path: 'scoreboard', Component: AdminContestScoreboard, apiRoute: 'GET /admin/contests/:contestID/scoreboard' },
  { path: 'teams', Component: AdminContestTeams, apiRoute: 'GET /admin/contests/:contestID/teams' },
  { path: 'teams/:teamId/details', Component: TeamDetails, apiRoute: 'GET /admin/contests/:contestID/teams' },
  { path: 'notices', Component: AdminContestNotices, apiRoute: 'GET /admin/contests/:contestID/notices' },
  { path: 'images', Component: AdminContestImagesPull, apiRoute: 'GET /admin/contests/:contestID/images' },
  { path: 'victims', Component: ContestContainers, apiRoute: 'GET /admin/contests/:contestID/victims' },
  { path: 'cheats', Component: AdminContestCheats, apiRoute: 'GET /admin/contests/:contestID/cheats' },
  { path: 'generators', Component: AdminContestGenerators, apiRoute: 'GET /admin/contests/:contestID/generators' },
];

export function AdminContestRoutes() {
  return (
    <Route
      path="/admin/contests/:id"
      element={withGuard(<AdminContestsLayout />, AdminRoute, { apiRoute: 'GET /admin/contests/:contestID' })}
    >
      <Route index element={withSuspense(AdminContestDetail)} />
      <Route path="settings" element={withSuspense(AdminContestSettings)} />
      {ADMIN_CONTEST_ROUTES.map(({ path, Component, apiRoute }) => (
        <Route key={path} path={path} element={withGuard(withSuspense(Component), AdminRoute, { apiRoute })} />
      ))}
    </Route>
  );
}
