import { useMemo, useState, useEffect } from 'react';
import { Outlet, useParams, useNavigate, useLocation } from 'react-router-dom';
import BaseLayout from './BaseLayout';
import { useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import { getContestInfo } from '../../../api/contest';
import { isContestEnded } from '../../../config/contest';
import { useBranding } from '../../../hooks/useBranding';

function ContestLayout() {
  const { contestId } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const user = useSelector((state) => state.user);
  const { t } = useTranslation();
  const { siteName } = useBranding();
  const [contestEnded, setContestEnded] = useState(false);

  useEffect(() => {
    getContestInfo(contestId).then((res) => {
      if (res.code === 200) {
        const contest = res.data;
        setContestEnded(isContestEnded(contest.start, contest.duration));
      }
    });
  }, [contestId]);

  const contestTabs = useMemo(() => {
    const tabs = [
      { id: 'overview', label: t('nav.overview') },
      { id: 'challenges', label: t('nav.challenges') },
      { id: 'scoreboard', label: t('nav.scoreboard') },
      { id: 'team', label: t('nav.team') },
      { id: 'notice', label: t('nav.notice') },
    ];

    if (!contestEnded) {
      tabs.splice(4, 0, { id: 'writeup', label: t('nav.writeup') });
    }

    return tabs;
  }, [t, contestEnded]);

  const activeTab = useMemo(() => {
    const currentTab = location.pathname.split('/').pop();
    return contestTabs.some((tab) => tab.id === currentTab) ? currentTab : 'overview';
  }, [location.pathname, contestTabs]);

  const handleTabChange = (tabId) => {
    if (tabId === 'overview') {
      navigate(`/contests/${contestId}`);
    } else {
      navigate(`/contests/${contestId}/${tabId}`);
    }
  };

  const handleLogoClick = () => {
    location.pathname !== `/contests/${contestId}` ? navigate(`/contests/${contestId}`) : navigate('/');
  };

  const handlePictureClick = () => {
    if (user.user) {
      navigate('/settings');
    } else {
      navigate('/login');
    }
  };

  return (
    <BaseLayout
      tabs={contestTabs}
      activeTab={activeTab}
      onTabChange={handleTabChange}
      onLogoClick={handleLogoClick}
      onPictureClick={handlePictureClick}
      logo={siteName}
      pictureSrc={user.user?.picture}
      userName={user.user?.name}
    >
      <div className="w-full">
        <Outlet />
      </div>
    </BaseLayout>
  );
}

export default ContestLayout;
