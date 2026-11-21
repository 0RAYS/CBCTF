import { lazy } from 'react';
import { Route } from 'react-router-dom';
import { UserRoute } from '@/routes/AuthRoute';
import { withGuard, withSuspense } from '@/routes/routeUtils';
const ContestLayout = lazy(() => import('@/components/features/layouts/ContestLayout'));

const GameDetailPage = lazy(() => import('@/pages/user/GameDetailPage'));
const GameChallengesPage = lazy(() => import('@/pages/user/GameChallengesPage'));
const GameLeaderboardPage = lazy(() => import('@/pages/user/GameLeaderboardPage'));
const GameTeamPage = lazy(() => import('@/pages/user/GameTeamPage'));
const GameNoticePage = lazy(() => import('@/pages/user/GameNoticePage'));
const GameWriteupPage = lazy(() => import('@/pages/user/GameWriteupPage'));

export function ContestRoutes() {
  return (
    <Route path="/contests/:contestId" element={withGuard(withSuspense(ContestLayout), UserRoute)}>
      <Route index element={withSuspense(GameDetailPage)} />
      <Route path="challenges" element={withSuspense(GameChallengesPage)} />
      <Route path="leaderboard" element={withSuspense(GameLeaderboardPage)} />
      <Route path="team" element={withSuspense(GameTeamPage)} />
      <Route path="notice" element={withSuspense(GameNoticePage)} />
      <Route path="writeup" element={withSuspense(GameWriteupPage)} />
    </Route>
  );
}
