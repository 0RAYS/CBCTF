import { useEffect } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useSelector } from 'react-redux';
import Loading from '../components/common/Loading';
import { useTranslation } from 'react-i18next';

export const UserRoute = ({ children, requiresAuth = true }) => {
  const { isAuthenticated, hasUserAccess, loading } = useSelector((state) => state.user);
  const location = useLocation();
  const { t } = useTranslation();

  useEffect(() => {
    document.title = `CBCTF - ${location.pathname.split('/').pop() || 'Home'}`;
  }, [location]);

  if (loading) {
    return <div className="flex items-center justify-center h-screen">{t('common.loading')}</div>;
  }

  if (requiresAuth && !isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // routes 尚未加载完成时放行，加载完成后若无权则拦截
  if (requiresAuth && isAuthenticated && !loading && !hasUserAccess) {
    return <Navigate to="/" replace />;
  }

  return children;
};

export const AdminRoute = ({ children, requiresAuth = true, apiRoute = null }) => {
  const { isAuthenticated, hasAdminAccess, loading, routes } = useSelector((state) => state.user);
  const location = useLocation();

  useEffect(() => {
    document.title = `CBCTF - ${location.pathname.split('/').pop() || 'Home'}`;
  }, [location]);

  if (loading) {
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
