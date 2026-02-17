import { useEffect } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useSelector } from 'react-redux';
import Loading from '../components/common/Loading';
import { useTranslation } from 'react-i18next';
export const UserRoute = ({ children, requiresAuth = true }) => {
  const { isAuthenticated, isAdmin, loading } = useSelector((state) => state.user);
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

  if (requiresAuth && isAuthenticated && isAdmin) {
    return <Navigate to="/admin/dashboard" state={{ from: location }} replace />;
  }

  return children;
};

export const AdminRoute = ({ children, requiresAuth = true }) => {
  const { isAuthenticated, isAdmin, loading } = useSelector((state) => state.user);

  const location = useLocation();
  useEffect(() => {
    document.title = `CBCTF - ${location.pathname.split('/').pop() || 'Home'}`;
  }, [location]);

  if (loading) {
    return <Loading />;
  }

  if (requiresAuth && !isAuthenticated) {
    return <Navigate to="/admin/login" state={{ from: location }} replace />;
  }

  if (requiresAuth && isAuthenticated && !isAdmin) {
    return <Navigate to="/" replace />;
  }

  return children;
};
