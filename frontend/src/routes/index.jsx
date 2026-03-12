import { Routes, Route } from 'react-router-dom';
import { Suspense, lazy } from 'react';
import { AdminRoute, UserRoute } from './AuthRoute';
import Loading from '../components/common/Loading';
import ErrorBoundary from '../components/common/ErrorBoundary';

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
const AdminContestImagesPull = lazy(() => import('../pages/admin/contests/images-pull.jsx'));
const ContestContainers = lazy(() => import('../pages/admin/contests/containers'));
const AdminContestCheats = lazy(() => import('../pages/admin/contests/cheats'));
const AdminContestGenerators = lazy(() => import('../pages/admin/contests/generators.jsx'));
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
    <ErrorBoundary>
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
                <Settings />
              </Suspense>
            }
          />
          <Route
            path="dashboard"
            element={
              <AdminRoute apiRoute="GET /admin/system/status">
                <Suspense fallback={<Loading />}>
                  <Dashboard />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="contests"
            element={
              <AdminRoute apiRoute="GET /admin/contests">
                <Suspense fallback={<Loading />}>
                  <ContestsManagement />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="rbac"
            element={
              <AdminRoute apiRoute="GET /admin/roles">
                <Suspense fallback={<Loading />}>
                  <RbacManagement />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="challenges"
            element={
              <AdminRoute apiRoute="GET /admin/challenges">
                <Suspense fallback={<Loading />}>
                  <ChallengesManagement />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="oauth"
            element={
              <AdminRoute apiRoute="GET /admin/oauth">
                <Suspense fallback={<Loading />}>
                  <OAuthProvidersManagement />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="smtp"
            element={
              <AdminRoute apiRoute="GET /admin/smtp">
                <Suspense fallback={<Loading />}>
                  <SmtpManagement />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="webhook"
            element={
              <AdminRoute apiRoute="GET /admin/webhook">
                <Suspense fallback={<Loading />}>
                  <WebhookManagement />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="files"
            element={
              <AdminRoute apiRoute="GET /admin/files">
                <Suspense fallback={<Loading />}>
                  <FilesManagement />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="system"
            element={
              <AdminRoute apiRoute="GET /admin/system/config">
                <Suspense fallback={<Loading />}>
                  <SystemSettings />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="logs"
            element={
              <AdminRoute apiRoute="GET /admin/logs">
                <Suspense fallback={<Loading />}>
                  <AdminLogs />
                </Suspense>
              </AdminRoute>
            }
          />
        </Route>

        {/* 管理端比赛详情路由 */}
        <Route
          path="/admin/contests/:id"
          element={
            <AdminRoute apiRoute="GET /admin/contests/:contestID">
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
              <AdminRoute apiRoute="GET /admin/contests/:contestID/challenges">
                <Suspense fallback={<Loading />}>
                  <AdminContestChallenges />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="scoreboard"
            element={
              <AdminRoute apiRoute="GET /admin/contests/:contestID/scoreboard">
                <Suspense fallback={<Loading />}>
                  <AdminContestScoreboard />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="teams"
            element={
              <AdminRoute apiRoute="GET /admin/contests/:contestID/teams">
                <Suspense fallback={<Loading />}>
                  <AdminContestTeams />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="teams/:teamId/details"
            element={
              <AdminRoute apiRoute="GET /admin/contests/:contestID/teams">
                <Suspense fallback={<Loading />}>
                  <TeamDetails />
                </Suspense>
              </AdminRoute>
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
              <AdminRoute apiRoute="GET /admin/contests/:contestID/notices">
                <Suspense fallback={<Loading />}>
                  <AdminContestNotices />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="images"
            element={
              <AdminRoute apiRoute="GET /admin/contests/:contestID/images">
                <Suspense fallback={<Loading />}>
                  <AdminContestImagesPull />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="containers"
            element={
              <AdminRoute apiRoute="GET /admin/contests/:contestID/victims">
                <Suspense fallback={<Loading />}>
                  <ContestContainers />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="cheats"
            element={
              <AdminRoute apiRoute="GET /admin/contests/:contestID/cheats">
                <Suspense fallback={<Loading />}>
                  <AdminContestCheats />
                </Suspense>
              </AdminRoute>
            }
          />
          <Route
            path="generators"
            element={
              <AdminRoute apiRoute="GET /admin/contests/:contestID/generators">
                <Suspense fallback={<Loading />}>
                  <AdminContestGenerators />
                </Suspense>
              </AdminRoute>
            }
          />
        </Route>
      </Routes>
    </ErrorBoundary>
  );
};

export default AppRoutes;
