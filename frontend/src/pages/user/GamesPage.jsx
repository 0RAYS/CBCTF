import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from '../../utils/toast';
import GameList from '../../components/features/Games/GameList';
import TeamJoinModal from '../../components/features/CTFGame/Team/TeamJoinModal';
import { getContestList } from '../../api/contest';
import { useSelector } from 'react-redux';
import { getTeamInfo, createTeam, joinTeam } from '../../api/game/team';
import Loading from '../../components/common/Loading';
import EmptyState from '../../components/common/EmptyState';
import { useTranslation } from 'react-i18next';
import { DEFAULT_CONTEST_IMAGE, getContestStatus, getContestTimeRange } from '../../config/contest';
import { IconTrophy, IconFlag, IconMedal } from '@tabler/icons-react';

const WELCOME_KEY = 'cbctf-welcome-dismissed';

// 转换比赛数据为组件需要的格式
const transformContestData = (contests) => {
  return contests.map((contest) => {
    const { startTime, endTime } = getContestTimeRange(contest.start, contest.duration);

    return {
      id: contest.id,
      title: `${contest.prefix} ${contest.name}`,
      description: contest.description,
      status: getContestStatus(contest.start, contest.duration),
      startTime,
      endTime,
      image: contest.picture || DEFAULT_CONTEST_IMAGE,
      teamSize: contest.size,
      teamsCount: contest.teams,
      usersCount: contest.users,
      noticesCount: contest.notices,
      isBlood: contest.blood,
      isHidden: contest.hidden,
    };
  });
};

function WelcomeBanner({ onDismiss }) {
  const { t } = useTranslation();

  const steps = [
    { num: '01', Icon: IconTrophy, label: t('game.welcome.step1'), desc: t('game.welcome.step1Desc') },
    { num: '02', Icon: IconFlag, label: t('game.welcome.step2'), desc: t('game.welcome.step2Desc') },
    { num: '03', Icon: IconMedal, label: t('game.welcome.step3'), desc: t('game.welcome.step3Desc') },
  ];

  return (
    <div className="border border-geek-400/30 bg-geek-400/5 rounded-md p-5 mb-6">
      <div className="flex items-start justify-between mb-5">
        <div>
          <h2 className="text-base font-mono text-neutral-50 mb-1">{t('game.welcome.title')}</h2>
          <p className="text-sm text-neutral-400">{t('game.welcome.subtitle')}</p>
        </div>
        <button
          onClick={onDismiss}
          className="text-neutral-400 hover:text-neutral-50 transition-colors text-sm font-mono ml-4 shrink-0 flex items-center gap-1"
          aria-label={t('game.welcome.dismiss')}
        >
          {t('game.welcome.dismiss')} ✕
        </button>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-5">
        {steps.map((step) => (
          <div key={step.num} className="flex items-start gap-3">
            <step.Icon className="w-5 h-5 shrink-0 mt-0.5 text-geek-400" />
            <div>
              <div className="text-sm font-mono text-neutral-50 mb-0.5">{step.label}</div>
              <div className="text-sm text-neutral-400">{step.desc}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

function GamesPage() {
  const [games, setGames] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentGameId, setCurrentGameId] = useState(null);
  const [showWelcome, setShowWelcome] = useState(() => !localStorage.getItem(WELCOME_KEY));
  const navigate = useNavigate();
  const user = useSelector((state) => state.user);
  const { t } = useTranslation();

  useEffect(() => {
    const fetchGames = async () => {
      try {
        const res = await getContestList();
        if (res.code === 200) {
          const transformedGames = transformContestData(res.data.contests);
          setGames(transformedGames);
        }
      } catch (error) {
        toast.danger({ title: t('toast.game.fetchListFailed'), description: error.message });
        setGames([]);
      } finally {
        setLoading(false);
      }
    };
    fetchGames();
  }, []);

  const handleDismissWelcome = () => {
    localStorage.setItem(WELCOME_KEY, '1');
    setShowWelcome(false);
  };

  const handleGameAction = async (gameId, action) => {
    if (!user.user) {
      navigate('/login');
      return;
    }

    const game = games.find((g) => g.id === gameId);
    if (!game) return;

    try {
      const checkResponse = await getTeamInfo(gameId);
      if (checkResponse.code === 200) {
        navigate(`/contests/${gameId}`);
      } else {
        switch (action) {
          case 'join':
            setCurrentGameId(gameId);
            setIsModalOpen(true);
            break;
          default:
            navigate(`/games`);
            break;
        }
      }
    } catch (error) {
      toast.danger({
        title: t('toast.team.checkStatusFailed'),
        description: error.message,
      });
    }
  };

  const handleCreateTeam = async (formData) => {
    try {
      const response = await createTeam(currentGameId, {
        name: formData.teamName,
        description: formData.description,
        captcha: formData.contestCode,
      });
      if (response.code === 200) {
        toast.success({ description: t('toast.team.createSuccess') });
        navigate(`/contests/${currentGameId}`);
      }
    } catch (error) {
      toast.danger({ title: t('toast.team.createFailed'), description: error.message });
    }
  };

  const handleJoinTeam = async (formData) => {
    try {
      const response = await joinTeam(currentGameId, {
        name: formData.teamName,
        captcha: formData.teamCode,
      });
      if (response.code === 200) {
        toast.success({ description: t('toast.team.joinSuccess') });
        navigate(`/contests/${currentGameId}`);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('toast.team.joinFailed') });
    }
  };

  if (loading) {
    return <Loading />;
  }

  return (
    <div>
      {showWelcome && <WelcomeBanner onDismiss={handleDismissWelcome} />}

      {games.length === 0 ? (
        <div className="py-16">
          <EmptyState title={t('game.noGames')} description={t('game.noGamesDescription')} />
        </div>
      ) : (
        <GameList
          games={games}
          onGameAction={handleGameAction}
          user={user}
        />
      )}
      <TeamJoinModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onCreateTeam={handleCreateTeam}
        onJoinTeam={handleJoinTeam}
      />
    </div>
  );
}

export default GamesPage;
