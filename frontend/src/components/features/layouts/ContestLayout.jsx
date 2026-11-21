import { useEffect, useMemo, useState } from 'react';
import { Outlet, useLocation, useNavigate, useParams } from 'react-router-dom';
import BaseLayout from '@/components/features/layouts/BaseLayout';
import { useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import { getContestInfo } from '@/api/contest';
import { getContestStatus } from '@/config/contest';
import { useBranding } from '@/hooks/useBranding';

function ContestLayout() {
  const { contestId } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const user = useSelector((state) => state.user);
  const { t } = useTranslation();
  const { siteName } = useBranding();
  const [contestStatus, setContestStatus] = useState(null);
  const [contestPicture, setContestPicture] = useState('');
  const [leaderboardHidden, setLeaderboardHidden] = useState(false);

  useEffect(() => {
    setContestPicture('');
    getContestInfo(contestId).then((res) => {
      if (res.code === 200) {
        const contest = res.data;
        setContestStatus(contest.paused ? 'paused' : getContestStatus(contest.start, contest.duration));
        setContestPicture(contest.picture || '');
        setLeaderboardHidden(!!contest.leaderboard_hidden);
      }
    });
  }, [contestId]);

  useEffect(() => {
    // 黑灯期间排行榜页面不可访问，退回比赛首页。
    if (leaderboardHidden && location.pathname.endsWith('/leaderboard')) {
      navigate(`/contests/${contestId}`);
    }
  }, [leaderboardHidden, location.pathname, contestId, navigate]);

  const contestTabs = useMemo(() => {
    const tabs = [
      { id: 'overview', label: t('nav.overview') },
      { id: 'challenges', label: t('nav.challenges') },
      { id: 'leaderboard', label: t('nav.leaderboard') },
      { id: 'team', label: t('nav.team') },
      { id: 'notice', label: t('nav.notice') },
    ];

    if (contestStatus && contestStatus !== 'upcoming' && contestStatus !== 'paused') {
      tabs.splice(4, 0, { id: 'writeup', label: t('nav.writeup') });
    }

    // 黑灯期间选手看不到排行榜入口。
    if (leaderboardHidden) {
      return tabs.filter((tab) => tab.id !== 'leaderboard');
    }

    return tabs;
  }, [t, contestStatus, leaderboardHidden]);

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
      backgroundImage={contestPicture}
    >
      <div className="w-full">
        <Outlet />
      </div>
    </BaseLayout>
  );
}

export default ContestLayout;
