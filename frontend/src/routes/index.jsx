import { Routes, useLocation } from 'react-router-dom';
import ErrorBoundary from '@/components/common/ErrorBoundary';
import { MainRoutes } from '@/routes/mainRoutes';
import { ContestRoutes } from '@/routes/contestRoutes';
import { AdminRoutes } from '@/routes/adminRoutes';
import { AdminContestRoutes } from '@/routes/adminContestRoutes';

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
