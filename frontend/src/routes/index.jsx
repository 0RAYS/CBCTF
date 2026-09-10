import { Routes, useLocation } from 'react-router-dom';
import ErrorBoundary from '../components/common/ErrorBoundary';
import { MainRoutes } from './mainRoutes';
import { ContestRoutes } from './contestRoutes';
import { AdminRoutes } from './adminRoutes';
import { AdminContestRoutes } from './adminContestRoutes';

const AppRoutes = () => {
  const { pathname } = useLocation();
  return (
    <ErrorBoundary resetKey={pathname}>
      <Routes>
        {MainRoutes()}
        {ContestRoutes()}
        {AdminRoutes()}
        {AdminContestRoutes()}
      </Routes>
    </ErrorBoundary>
  );
};

export default AppRoutes;
