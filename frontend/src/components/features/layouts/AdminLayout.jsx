import { Outlet, useNavigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import AdminShellLayout from './AdminShellLayout';
import { getAdminNavSections } from '../../../config/adminNavigation';

function AdminLayout() {
  const navigate = useNavigate();
  const user = useSelector((state) => state.user);
  const { t } = useTranslation();

  const sections = useMemo(() => getAdminNavSections(t), [t]);

  // Handle logo click
  const handleLogoClick = () => {
    navigate('/');
  };

  // Handle picture click
  const handlePictureClick = () => {
    if (user.user) {
      navigate('/admin/settings');
    } else {
      navigate('/admin/login');
    }
  };

  return (
    <AdminShellLayout
      sections={sections}
      onNavigate={navigate}
      onLogoClick={handleLogoClick}
      onAvatarClick={handlePictureClick}
      user={user}
      logo={t('branding.admin')}
    >
      <Outlet />
    </AdminShellLayout>
  );
}

export default AdminLayout;
