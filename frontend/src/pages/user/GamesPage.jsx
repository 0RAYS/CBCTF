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
import Button from '../../components/common/Button';
import { useTranslation } from 'react-i18next';
import { DEFAULT_CONTEST_IMAGE, getContestStatus, getContestTimeRange } from '../../config/contest';

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

function GamesPage() {
  const [games, setGames] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [retry, setRetry] = useState(0);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentGameId, setCurrentGameId] = useState(null);
  const navigate = useNavigate();
  const user = useSelector((state) => state.user);
  const { t } = useTranslation();

  useEffect(() => {
    let active = true;
    const fetchGames = async () => {
      setLoading(true);
      setError(null);
      try {
        const res = await getContestList();
        if (!active) return;
        if (res.code !== 200) throw new Error(res.msg || t('toast.game.fetchListFailed'));
        setGames(transformContestData(res.data?.contests || []));
      } catch (error) {
        if (active) setError(error.message || t('errors.requestFailed'));
      } finally {
        if (active) setLoading(false);
      }
    };
    fetchGames();
    return () => {
      active = false;
    };
  }, [retry, t]);

  const handleGameAction = async (gameId, action) => {
    if (!user.user) {
      navigate('/login');
      return;
    }

    const game = games.find((g) => g.id === gameId);
    if (!game) return;

    try {
      const checkResponse = await getTeamInfo(gameId, { noToast: true });
      if (checkResponse.code === 200) {
        navigate(`/contests/${gameId}`);
      } else {
        switch (action) {
          case 'join':
            setCurrentGameId(gameId);
            setIsModalOpen(true);
            break;
          default:
            if (checkResponse.code === 404) {
              toast.warning({ title: t('toast.team.notFound'), description: t('toast.team.notJoinContest') });
            }
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
    const response = await createTeam(currentGameId, {
      name: formData.teamName,
      description: formData.description,
      captcha: formData.contestCode,
    });
    if (response.code !== 200) throw new Error(response.msg || t('toast.team.createFailed'));
    toast.success({ description: t('toast.team.createSuccess') });
    navigate(`/contests/${currentGameId}`);
    return true;
  };

  const handleJoinTeam = async (formData) => {
    const response = await joinTeam(currentGameId, {
      name: formData.teamName,
      captcha: formData.teamCode,
    });
    if (response.code !== 200) throw new Error(response.msg || t('toast.team.joinFailed'));
    toast.success({ description: t('toast.team.joinSuccess') });
    navigate(`/contests/${currentGameId}`);
    return true;
  };

  if (loading) {
    return <Loading />;
  }

  return (
    <div>
      {error ? (
        <div role="alert" className="py-16">
          <EmptyState
            title={t('toast.game.fetchListFailed')}
            description={error}
            action={<Button onClick={() => setRetry((value) => value + 1)}>{t('common.refresh')}</Button>}
          />
        </div>
      ) : games.length === 0 ? (
        <div className="py-16">
          <EmptyState title={t('game.noGames')} description={t('game.noGamesDescription')} />
        </div>
      ) : (
        <GameList games={games} onGameAction={handleGameAction} user={user} />
      )}
      <TeamJoinModal
        key={`${currentGameId}-${isModalOpen}`}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onCreateTeam={handleCreateTeam}
        onJoinTeam={handleJoinTeam}
      />
    </div>
  );
}

export default GamesPage;
