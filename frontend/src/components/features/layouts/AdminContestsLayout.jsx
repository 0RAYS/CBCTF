import { useMemo } from 'react';
import { Outlet, useNavigate, useParams } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import AdminShellLayout from './AdminShellLayout';
import { getAdminContestNavSections } from '../../../config/adminNavigation';

function AdminContestsLayout() {
  const { id } = useParams();
  const navigate = useNavigate();
  const user = useSelector((state) => state.user);
  const { t } = useTranslation();

  const sections = useMemo(() => getAdminContestNavSections(t, id), [t, id]);

  const handleLogoClick = () => {
    navigate('/admin/dashboard');
  };

  const handlePictureClick = () => {
    navigate('/admin/settings');
  };

  return (
    <AdminShellLayout
      sections={sections}
      onNavigate={navigate}
      onLogoClick={handleLogoClick}
      onAvatarClick={handlePictureClick}
      user={user}
      logo={t('branding.admin')}
      subtitle={id ? `${t('admin.contest')} #${id}` : t('admin.contest')}
    >
      <div className="w-full">
        <Outlet />
      </div>
    </AdminShellLayout>
  );
}

export default AdminContestsLayout;
