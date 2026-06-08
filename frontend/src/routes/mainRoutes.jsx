import { lazy } from 'react';
import { Route } from 'react-router-dom';
import { UserRoute } from './AuthRoute';
import { withGuard, withSuspense } from './routeUtils';
import MainLayout from '../components/features/layouts/MainLayout';

const Login = lazy(() => import('../pages/Login'));
const OAuthCallback = lazy(() => import('../pages/OAuthCallback'));
const Home = lazy(() => import('../pages/Home'));
const Settings = lazy(() => import('../pages/user/Settings'));
const GamesPage = lazy(() => import('../pages/user/GamesPage'));
const ResetPassword = lazy(() => import('../pages/user/ResetPassword'));
const VerifyEmail = lazy(() => import('../pages/user/VerifyEmail'));

export function MainRoutes() {
  return (
    <Route path="/" element={<MainLayout />}>
      <Route index element={withSuspense(Home)} />
      <Route path="settings" element={withGuard(withSuspense(Settings), UserRoute)} />
      <Route path="games" element={withSuspense(GamesPage)} />
      <Route path="login" element={withSuspense(Login)} />
      <Route path="oauth/callback" element={withSuspense(OAuthCallback)} />
      <Route path="reset-password" element={withSuspense(ResetPassword)} />
      <Route path="verify" element={withSuspense(VerifyEmail)} />
    </Route>
  );
}
