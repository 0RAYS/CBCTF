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
const CronJobsManagement = lazy(() => import('../pages/admin/CronJobs'));
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
const AdminVictims = lazy(() => import('../pages/admin/Victims'));
const AdminGenerators = lazy(() => import('../pages/admin/Generators'));
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

const withSuspense = (Component) => (
  <Suspense fallback={<Loading />}>
    <Component />
  </Suspense>
);

const withGuard = (element, Guard, guardProps = {}) => <Guard {...guardProps}>{element}</Guard>;

const AppRoutes = () => {
  return (
    <ErrorBoundary>
      <Routes>
        {/* 主布局路由 */}
        <Route path="/" element={<MainLayout />}>
          <Route index element={withSuspense(Home)} />
          <Route path="settings" element={withGuard(withSuspense(Settings), UserRoute)} />
          <Route path="games" element={withSuspense(GamesPage)} />
          <Route path="login" element={withSuspense(Login)} />
          <Route path="oauth/callback" element={withSuspense(OAuthCallback)} />
          <Route path="support" element={withSuspense(TechStackPage)} />
          <Route path="contact" element={withSuspense(ContactPage)} />
        </Route>

        {/* 比赛详情路由 */}
        <Route path="/contests/:contestId" element={withGuard(<ContestLayout />, UserRoute)}>
          <Route index element={withSuspense(GameDetailPage)} />
          <Route path="challenges" element={withSuspense(GameChallengesPage)} />
          <Route path="scoreboard" element={withSuspense(GameScoreBoardPage)} />
          <Route path="team" element={withSuspense(GameTeamPage)} />
          <Route path="notice" element={withSuspense(GameNoticePage)} />
          <Route path="writeup" element={withSuspense(GameWriteupPage)} />
        </Route>

        {/* 管理员路由 */}
        <Route path="/admin" element={withGuard(<AdminLayout />, AdminRoute)}>
          <Route path="settings" element={withSuspense(Settings)} />
          <Route
            path="dashboard"
            element={withGuard(withSuspense(Dashboard), AdminRoute, { apiRoute: 'GET /admin/system/status' })}
          />
          <Route
            path="contests"
            element={withGuard(withSuspense(ContestsManagement), AdminRoute, { apiRoute: 'GET /admin/contests' })}
          />
          <Route
            path="rbac"
            element={withGuard(withSuspense(RbacManagement), AdminRoute, { apiRoute: 'GET /admin/roles' })}
          />
          <Route
            path="challenges"
            element={withGuard(withSuspense(ChallengesManagement), AdminRoute, { apiRoute: 'GET /admin/challenges' })}
          />
          <Route
            path="oauth"
            element={withGuard(withSuspense(OAuthProvidersManagement), AdminRoute, { apiRoute: 'GET /admin/oauth' })}
          />
          <Route
            path="smtp"
            element={withGuard(withSuspense(SmtpManagement), AdminRoute, { apiRoute: 'GET /admin/smtp' })}
          />
          <Route
            path="cronjobs"
            element={withGuard(withSuspense(CronJobsManagement), AdminRoute, { apiRoute: 'GET /admin/cronjobs' })}
          />
          <Route
            path="webhook"
            element={withGuard(withSuspense(WebhookManagement), AdminRoute, { apiRoute: 'GET /admin/webhook' })}
          />
          <Route
            path="files"
            element={withGuard(withSuspense(FilesManagement), AdminRoute, { apiRoute: 'GET /admin/files' })}
          />
          <Route
            path="system"
            element={withGuard(withSuspense(SystemSettings), AdminRoute, { apiRoute: 'GET /admin/system/config' })}
          />
          <Route
            path="logs"
            element={withGuard(withSuspense(AdminLogs), AdminRoute, { apiRoute: 'GET /admin/logs' })}
          />
          <Route
            path="victims"
            element={withGuard(withSuspense(AdminVictims), AdminRoute, { apiRoute: 'GET /admin/victims' })}
          />
          <Route
            path="generators"
            element={withGuard(withSuspense(AdminGenerators), AdminRoute, { apiRoute: 'GET /admin/generators' })}
          />
        </Route>

        {/* 管理端比赛详情路由 */}
        <Route
          path="/admin/contests/:id"
          element={withGuard(<AdminContestsLayout />, AdminRoute, { apiRoute: 'GET /admin/contests/:contestID' })}
        >
          <Route index element={withSuspense(AdminContestDetail)} />
          <Route
            path="challenges"
            element={withGuard(withSuspense(AdminContestChallenges), AdminRoute, {
              apiRoute: 'GET /admin/contests/:contestID/challenges',
            })}
          />
          <Route
            path="scoreboard"
            element={withGuard(withSuspense(AdminContestScoreboard), AdminRoute, {
              apiRoute: 'GET /admin/contests/:contestID/scoreboard',
            })}
          />
          <Route
            path="teams"
            element={withGuard(withSuspense(AdminContestTeams), AdminRoute, {
              apiRoute: 'GET /admin/contests/:contestID/teams',
            })}
          />
          <Route
            path="teams/:teamId/details"
            element={withGuard(withSuspense(TeamDetails), AdminRoute, {
              apiRoute: 'GET /admin/contests/:contestID/teams',
            })}
          />
          <Route path="settings" element={withSuspense(AdminContestSettings)} />
          <Route
            path="notices"
            element={withGuard(withSuspense(AdminContestNotices), AdminRoute, {
              apiRoute: 'GET /admin/contests/:contestID/notices',
            })}
          />
          <Route
            path="images"
            element={withGuard(withSuspense(AdminContestImagesPull), AdminRoute, {
              apiRoute: 'GET /admin/contests/:contestID/images',
            })}
          />
          <Route
            path="victims"
            element={withGuard(withSuspense(ContestContainers), AdminRoute, {
              apiRoute: 'GET /admin/contests/:contestID/victims',
            })}
          />
          <Route
            path="cheats"
            element={withGuard(withSuspense(AdminContestCheats), AdminRoute, {
              apiRoute: 'GET /admin/contests/:contestID/cheats',
            })}
          />
          <Route
            path="generators"
            element={withGuard(withSuspense(AdminContestGenerators), AdminRoute, {
              apiRoute: 'GET /admin/contests/:contestID/generators',
            })}
          />
        </Route>
      </Routes>
    </ErrorBoundary>
  );
};

export default AppRoutes;
