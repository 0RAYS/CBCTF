import { Outlet, useNavigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import BaseLayout from './BaseLayout';

function AdminAuthLayout() {
  const user = useSelector((state) => state.user);
  const navigate = useNavigate();
  const { t } = useTranslation();

  return (
    <BaseLayout
      tabs={[]}
      activeTab=""
      onTabChange={() => {}}
      onLogoClick={() => navigate('/')}
      onPictureClick={() => navigate('/admin/login')}
      logo={t('branding.main')}
      pictureSrc={user.user?.picture}
      userName={user.user?.name}
    >
      <Outlet />
    </BaseLayout>
  );
}

export default AdminAuthLayout;
