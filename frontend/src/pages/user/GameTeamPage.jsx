import { useParams, useNavigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { useTranslation } from 'react-i18next';
import TeamSettings from '../../components/features/CTFGame/Team/TeamSettings';
import useTeamSettings from '../../components/features/CTFGame/Team/useTeamSettings';
import Loading from '../../components/common/Loading';

function GameTeamContent({ contestId, userId }) {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { team, loading, ...actions } = useTeamSettings(contestId, navigate);

  if (loading) return <Loading />;
  if (!team) {
    return (
      <div className="flex items-center justify-center min-h-[500px]">
        <span className="text-neutral-300">{t('game.team.empty')}</span>
      </div>
    );
  }
  return <TeamSettings team={team} isLeader={userId === team.captainId} {...actions} />;
}

export default function GameTeamPage() {
  const { contestId } = useParams();
  const userId = useSelector((state) => state.user.user?.id);
  return <GameTeamContent key={JSON.stringify([contestId, userId])} contestId={contestId} userId={userId} />;
}
