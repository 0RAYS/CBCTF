import HeroSection from '../components/features/Home/HeroSection';
import { lazy, startTransition, Suspense, useEffect, useState } from 'react';
import { getStats } from '../api/stats.js';
import { toast } from '../utils/toast.js';
import { useTranslation } from 'react-i18next';

const StatsSection = lazy(() => import('../components/features/Home/StatsSection'));
const ChallengeTypes = lazy(() => import('../components/features/Home/ChallengeTypes'));
const UpcomingContests = lazy(() => import('../components/features/Home/UpcomingContests'));
const LeaderboardPreview = lazy(() => import('../components/features/Home/LeaderboardPreview'));

const transformUpcomingData = (contests) => {
  return contests.map((contest) => {
    return {
      title: contest.name,
      date: new Date(contest.start).toLocaleDateString(),
      duration: `${contest.duration / 3600}h`,
      registrations: contest.users,
      teams: contest.teams,
      image: contest.picture,
    };
  });
};

function Home() {
  const { t } = useTranslation();
  const [stats, setStats] = useState([]);
  const [leaderboard, setLeaderboard] = useState([]);
  const [upcomingContests, setUpcomingContests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showDeferredSections, setShowDeferredSections] = useState(false);

  useEffect(() => {
    const showSections = () => {
      startTransition(() => setShowDeferredSections(true));
    };

    if ('requestIdleCallback' in window) {
      const idleId = window.requestIdleCallback(showSections, { timeout: 1200 });
      return () => window.cancelIdleCallback(idleId);
    }

    const timeoutId = window.setTimeout(showSections, 300);
    return () => window.clearTimeout(timeoutId);
  }, []);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const res = await getStats();
        if (res.code === 200) {
          setStats(res.data.stats);
          setLeaderboard(res.data.scoreboard);
          setUpcomingContests(transformUpcomingData(res.data.upcoming));
        }
      } catch (error) {
        toast.danger({ title: t('toast.home.fetchFailed'), description: error.message });
        setStats([]);
        setLeaderboard([]);
        setUpcomingContests([]);
      } finally {
        setLoading(false);
      }
    };
    fetchStats();
  }, []);

  return (
    <div className="-mt-10">
      <HeroSection />
      {showDeferredSections && (
        <Suspense fallback={null}>
          <StatsSection stats={stats} isLoading={loading} />
          <ChallengeTypes />
          <UpcomingContests contests={upcomingContests} isLoading={loading} />
          <LeaderboardPreview topUsers={leaderboard} isLoading={loading} />
        </Suspense>
      )}
    </div>
  );
}

export default Home;
