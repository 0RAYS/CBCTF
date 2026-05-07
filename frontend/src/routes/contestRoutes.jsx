import { lazy } from 'react';
import { Route } from 'react-router-dom';
import { UserRoute } from './AuthRoute';
import { withGuard, withSuspense } from './routeUtils';
import ContestLayout from '../components/features/layouts/ContestLayout';

const GameDetailPage = lazy(() => import('../pages/user/GameDetailPage'));
const GameChallengesPage = lazy(() => import('../pages/user/GameChallengesPage'));
const GameScoreBoardPage = lazy(() => import('../pages/user/GameScoreBoardPage'));
const GameTeamPage = lazy(() => import('../pages/user/GameTeamPage'));
const GameNoticePage = lazy(() => import('../pages/user/GameNoticePage'));
const GameWriteupPage = lazy(() => import('../pages/user/GameWriteupPage'));

export function ContestRoutes() {
  return (
    <Route path="/contests/:contestId" element={withGuard(<ContestLayout />, UserRoute)}>
      <Route index element={withSuspense(GameDetailPage)} />
      <Route path="challenges" element={withSuspense(GameChallengesPage)} />
      <Route path="scoreboard" element={withSuspense(GameScoreBoardPage)} />
      <Route path="team" element={withSuspense(GameTeamPage)} />
      <Route path="notice" element={withSuspense(GameNoticePage)} />
      <Route path="writeup" element={withSuspense(GameWriteupPage)} />
    </Route>
  );
}
