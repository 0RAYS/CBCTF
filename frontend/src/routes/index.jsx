import { Routes, Route } from 'react-router-dom';
import { Suspense, lazy } from 'react';
import { AdminRoute, UserRoute } from './AuthRoute';
import Loading from '../components/common/Loading';

// 直接导入重要的布局组件
import ContestLayout from '../components/features/layouts/ContestLayout';
import MainLayout from '../components/features/layouts/MainLayout';
import AdminLayout from '../components/features/layouts/AdminLayout';
import AdminContestsLayout from '../components/features/layouts/AdminContestsLayout';

// 懒加载页面组件
const Login = lazy(() => import('../pages/Login'));
const Dashboard = lazy(() => import('../pages/admin/Dashboard'));
const ContestsManagement = lazy(() => import('../pages/admin/Contests'));
const RbacManagement = lazy(() => import('../pages/admin/Rbac'));
const ChallengesManagement = lazy(() => import('../pages/admin/Challenges'));
const FilesManagement = lazy(() => import('../pages/admin/Files.jsx'));
const SystemSettings = lazy(() => import('../pages/admin/System'));
const AdminSettings = lazy(() => import('../pages/admin/Settings'));
const OAuthProvidersManagement = lazy(() => import('../pages/admin/OAuthProviders'));
const SmtpManagement = lazy(() => import('../pages/admin/Smtp'));
const WebhookManagement = lazy(() => import('../pages/admin/Webhook'));
const AdminLogs = lazy(() => import('../pages/admin/Logs'));
const OAuthCallback = lazy(() => import('../pages/OAuthCallback'));
const Home = lazy(() => import('../pages/Home'));

// 懒加载管理端比赛详情相关组件
const AdminContestDetail = lazy(() => import('../pages/admin/contests/index.jsx'));
const AdminContestChallenges = lazy(() => import('../pages/admin/contests/challenges'));
const AdminContestScoreboard = lazy(() => import('../pages/admin/contests/scoreboard'));
const AdminContestTeams = lazy(() => import('../pages/admin/contests/teams'));
const AdminContestSettings = lazy(() => import('../pages/admin/contests/settings'));
const TeamDetails = lazy(() => import('../pages/admin/contests/team-details'));
const AdminContestNotices = lazy(() => import('../pages/admin/contests/notice'));
const AdminContestImagesWarmup = lazy(() => import('../pages/admin/contests/images-warmup.jsx'));
const ContestContainers = lazy(() => import('../pages/admin/contests/containers'));
const AdminContestCheats = lazy(() => import('../pages/admin/contests/cheats'));
const Settings = lazy(() => import('../pages/user/Settings'));
const GamesPage = lazy(() => import('../pages/user/GamesPage'));
const GameDetailPage = lazy(() => import('../pages/user/GameDetailPage'));
const GameChallengesPage = lazy(() => import('../pages/user/GameChallengesPage'));
const GameScoreBoardPage = lazy(() => import('../pages/user/GameScoreBoardPage'));
const GameTeamPage = lazy(() => import('../pages/user/GameTeamPage'));
const GameNoticePage = lazy(() => import('../pages/user/GameNoticePage'));
const GameWriteupPage = lazy(() => import('../pages/user/GameWriteupPage'));
const TechStackPage = lazy(() => import('../pages/TechStackPage'));
const ContactPage = lazy(() => import('../pages/ContactPage'));

const AppRoutes = () => {
  return (
    <Routes>
      {/* 主布局路由 */}
      <Route path="/" element={<MainLayout />}>
        <Route
          index
          element={
            <Suspense fallback={<Loading />}>
              <Home />
            </Suspense>
          }
        />
        <Route
          path="settings"
          element={
            <UserRoute>
              <Suspense fallback={<Loading />}>
                <Settings />
              </Suspense>
            </UserRoute>
          }
        />
        <Route
          path="games"
          element={
            <Suspense fallback={<Loading />}>
              <GamesPage />
            </Suspense>
          }
        />
        <Route
          path="login"
          element={
            <Suspense fallback={<Loading />}>
              <Login />
            </Suspense>
          }
        />
        <Route
          path="oauth/callback"
          element={
            <Suspense fallback={<Loading />}>
              <OAuthCallback />
            </Suspense>
          }
        />
        <Route
          path="support"
          element={
            <Suspense fallback={<Loading />}>
              <TechStackPage />
            </Suspense>
          }
        />
        <Route
          path="contact"
          element={
            <Suspense fallback={<Loading />}>
              <ContactPage />
            </Suspense>
          }
        />
      </Route>

      {/* 比赛详情路由 */}
      <Route
        path="/contests/:contestId"
        element={
          <UserRoute>
            <ContestLayout />
          </UserRoute>
        }
      >
        <Route
          index
          element={
            <Suspense fallback={<Loading />}>
              <GameDetailPage />
            </Suspense>
          }
        />
        <Route
          path="challenges"
          element={
            <Suspense fallback={<Loading />}>
              <GameChallengesPage />
            </Suspense>
          }
        />
        <Route
          path="scoreboard"
          element={
            <Suspense fallback={<Loading />}>
              <GameScoreBoardPage />
            </Suspense>
          }
        />
        <Route
          path="team"
          element={
            <Suspense fallback={<Loading />}>
              <GameTeamPage />
            </Suspense>
          }
        />
        <Route
          path="notice"
          element={
            <Suspense fallback={<Loading />}>
              <GameNoticePage />
            </Suspense>
          }
        />
        <Route
          path="writeup"
          element={
            <Suspense fallback={<Loading />}>
              <GameWriteupPage />
            </Suspense>
          }
        />
      </Route>

      {/* 管理员路由 */}
      <Route
        path="/admin"
        element={
          <AdminRoute>
            <AdminLayout />
          </AdminRoute>
        }
      >
        <Route
          path="settings"
          element={
            <Suspense fallback={<Loading />}>
              <AdminSettings />
            </Suspense>
          }
        />
        <Route
          path="dashboard"
          element={
            <Suspense fallback={<Loading />}>
              <Dashboard />
            </Suspense>
          }
        />
        <Route
          path="contests"
          element={
            <Suspense fallback={<Loading />}>
              <ContestsManagement />
            </Suspense>
          }
        />
        <Route
          path="rbac"
          element={
            <Suspense fallback={<Loading />}>
              <RbacManagement />
            </Suspense>
          }
        />
        <Route
          path="challenges"
          element={
            <Suspense fallback={<Loading />}>
              <ChallengesManagement />
            </Suspense>
          }
        />
        <Route
          path="oauth"
          element={
            <Suspense fallback={<Loading />}>
              <OAuthProvidersManagement />
            </Suspense>
          }
        />
        <Route
          path="smtp"
          element={
            <Suspense fallback={<Loading />}>
              <SmtpManagement />
            </Suspense>
          }
        />
        <Route
          path="webhook"
          element={
            <Suspense fallback={<Loading />}>
              <WebhookManagement />
            </Suspense>
          }
        />
        <Route
          path="files"
          element={
            <Suspense fallback={<Loading />}>
              <FilesManagement />
            </Suspense>
          }
        />
        <Route
          path="system"
          element={
            <Suspense fallback={<Loading />}>
              <SystemSettings />
            </Suspense>
          }
        />
        <Route
          path="logs"
          element={
            <Suspense fallback={<Loading />}>
              <AdminLogs />
            </Suspense>
          }
        />
      </Route>

      {/* 管理端比赛详情路由 */}
      <Route
        path="/admin/contests/:id"
        element={
          <AdminRoute>
            <AdminContestsLayout />
          </AdminRoute>
        }
      >
        <Route
          index
          element={
            <Suspense fallback={<Loading />}>
              <AdminContestDetail />
            </Suspense>
          }
        />
        <Route
          path="challenges"
          element={
            <Suspense fallback={<Loading />}>
              <AdminContestChallenges />
            </Suspense>
          }
        />
        <Route
          path="scoreboard"
          element={
            <Suspense fallback={<Loading />}>
              <AdminContestScoreboard />
            </Suspense>
          }
        />
        <Route
          path="teams"
          element={
            <Suspense fallback={<Loading />}>
              <AdminContestTeams />
            </Suspense>
          }
        />
        <Route
          path="teams/:teamId/details"
          element={
            <Suspense fallback={<Loading />}>
              <TeamDetails />
            </Suspense>
          }
        />
        <Route
          path="settings"
          element={
            <Suspense fallback={<Loading />}>
              <AdminContestSettings />
            </Suspense>
          }
        />
        <Route
          path="notices"
          element={
            <Suspense fallback={<Loading />}>
              <AdminContestNotices />
            </Suspense>
          }
        />
        <Route
          path="images"
          element={
            <Suspense fallback={<Loading />}>
              <AdminContestImagesWarmup />
            </Suspense>
          }
        />
        <Route
          path="containers"
          element={
            <Suspense fallback={<Loading />}>
              <ContestContainers />
            </Suspense>
          }
        />
        <Route
          path="cheats"
          element={
            <Suspense fallback={<Loading />}>
              <AdminContestCheats />
            </Suspense>
          }
        />
      </Route>
    </Routes>
  );
};

export default AppRoutes;
