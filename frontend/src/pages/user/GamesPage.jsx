import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from '../../utils/toast';
import GameSlider from '../../components/features/Games/GameSlider';
import TeamJoinModal from '../../components/features/CTFGame/Team/TeamJoinModal';
import { getContestList } from '../../api/contest';
import { useSelector } from 'react-redux';
import { getTeamInfo, createTeam, joinTeam } from '../../api/game/team';
import Loading from '../../components/common/Loading';
import EmptyState from '../../components/common/EmptyState';
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
      // 额外信息
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
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentGameId, setCurrentGameId] = useState(null);
  const navigate = useNavigate();
  const user = useSelector((state) => state.user);
  const { t } = useTranslation();

  // 获取游戏列表
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
        setGames([]); // 出错时设置为空数组
      } finally {
        setLoading(false);
      }
    };

    fetchGames();
  }, []);

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
        // 用户已加入，直接跳转
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
      {games.length === 0 ? (
        <div className="py-24">
          <EmptyState title={t('game.noGames')} />
        </div>
      ) : (
        <GameSlider
          games={games}
          currentIndex={currentIndex}
          onIndexChange={setCurrentIndex}
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
