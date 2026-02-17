import HeroSection from '../components/features/Home/HeroSection';
import StatsSection from '../components/features/Home/StatsSection';
import ChallengeTypes from '../components/features/Home/ChallengeTypes';
import UpcomingContests from '../components/features/Home/UpcomingContests';
import LeaderboardPreview from '../components/features/Home/LeaderboardPreview';
import { useEffect, useState } from 'react';
import { getStats } from '../api/stats.js';
import { toast } from '../utils/toast.js';
import { useTranslation } from 'react-i18next';

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

  if (loading) {
    return (
      <>
        <HeroSection />
        <ChallengeTypes />
      </>
    );
  }

  return (
    <div className="-mt-10">
      <HeroSection />
      <StatsSection stats={stats} />
      <ChallengeTypes />
      <UpcomingContests contests={upcomingContests} />
      <LeaderboardPreview topUsers={leaderboard} />
    </div>
  );
}

export default Home;
