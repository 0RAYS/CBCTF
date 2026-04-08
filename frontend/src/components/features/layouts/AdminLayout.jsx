import { Outlet, useNavigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import AdminShellLayout from './AdminShellLayout';
import { getAdminNavSections } from '../../../config/adminNavigation';
import { useBranding } from '../../../hooks/useBranding';

function AdminLayout() {
  const navigate = useNavigate();
  const user = useSelector((state) => state.user);
  const routes = useSelector((state) => state.user.routes);
  const { t } = useTranslation();
  const { adminName } = useBranding();

  const sections = useMemo(() => getAdminNavSections(t, routes), [t, routes]);

  // Handle logo click
  const handleLogoClick = () => {
    navigate('/');
  };

  // Handle picture click
  const handlePictureClick = () => {
    if (user.user) {
      navigate('/admin/settings');
    } else {
      navigate('/login');
    }
  };

  return (
    <AdminShellLayout
      sections={sections}
      onNavigate={navigate}
      onLogoClick={handleLogoClick}
      onAvatarClick={handlePictureClick}
      user={user}
      logo={adminName}
    >
      <Outlet />
    </AdminShellLayout>
  );
}

export default AdminLayout;
