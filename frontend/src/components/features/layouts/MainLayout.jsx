import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import BaseLayout from '@/components/features/layouts/BaseLayout';
import { useSelector } from 'react-redux';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useBranding } from '@/hooks/useBranding';

function MainLayout() {
  const location = useLocation();
  const navigate = useNavigate();
  const user = useSelector((state) => state.user);
  const { t } = useTranslation();
  const { siteName } = useBranding();

  const mainTabs = useMemo(() => {
    const tabs = [{ id: 'GAMES', label: t('nav.games') }];
    if (user.isAuthenticated) tabs.push({ id: 'LEADERBOARD', label: t('nav.leaderboard') });
    return tabs;
  }, [t, user.isAuthenticated]);

  const activeTab =
    location.pathname === '/games' ? 'GAMES' : location.pathname === '/leaderboard' ? 'LEADERBOARD' : '';

  const handleTabChange = (tabId) => {
    navigate(tabId === 'LEADERBOARD' ? '/leaderboard' : '/games');
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
      activeTab={activeTab}
      onTabChange={handleTabChange}
      onLogoClick={handleLogoClick}
      onPictureClick={handlePictureClick}
      logo={siteName}
      pictureSrc={user.user?.picture}
      userName={user.user?.name}
    >
      <Outlet />
    </BaseLayout>
  );
}

export default MainLayout;
