import { Navigate, useLocation } from 'react-router-dom';
import { useSelector } from 'react-redux';
import Loading from '@/components/common/Loading';

export const UserRoute = ({ children, requiresAuth = true }) => {
  const { isAuthenticated, hasUserAccess, initialized } = useSelector((state) => state.user);
  const location = useLocation();

  if (!initialized) {
    return <Loading />;
  }

  if (requiresAuth && !isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  if (requiresAuth && isAuthenticated && !hasUserAccess) {
    return <Navigate to="/" replace />;
  }

  return children;
};

export const AdminRoute = ({ children, requiresAuth = true, apiRoute = null }) => {
  const { isAuthenticated, hasAdminAccess, initialized, routes } = useSelector((state) => state.user);
  const location = useLocation();

  if (!initialized) {
    return <Loading />;
  }

  if (requiresAuth && !isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  if (requiresAuth && isAuthenticated && !hasAdminAccess) {
    return <Navigate to="/" replace />;
  }

  if (apiRoute && routes.length > 0 && !routes.includes(apiRoute)) {
    return <Navigate to="/admin/dashboard" replace />;
  }

  return children;
};
