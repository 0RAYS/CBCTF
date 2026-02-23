import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import BaseLayout from './BaseLayout';
import { useSelector } from 'react-redux';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
function MainLayout() {
  const location = useLocation();
  const navigate = useNavigate();
  const user = useSelector((state) => state.user);
  const { t } = useTranslation();

  const mainTabs = useMemo(() => [{ id: 'GAMES', label: t('nav.games') }], [t]);

  // 根据路径设置激活状态
  const isGamesActive = location.pathname === '/games';

  const handleTabChange = () => {
    navigate('/games');
  };

  const handleLogoClick = () => {
    navigate('/');
  };

  const handlePictureClick = () => {
    if (user.user) {
      navigate(user.hasAdminAccess ? '/admin/settings' : '/settings');
    } else {
      navigate('/login');
    }
  };

  return (
    <BaseLayout
      tabs={mainTabs}
      activeTab={isGamesActive ? 'GAMES' : ''}
      onTabChange={handleTabChange}
      onLogoClick={handleLogoClick}
      onPictureClick={handlePictureClick}
      logo={t('branding.main')}
      pictureSrc={user.user?.picture}
      userName={user.user?.name}
    >
      <Outlet />
    </BaseLayout>
  );
}

export default MainLayout;
