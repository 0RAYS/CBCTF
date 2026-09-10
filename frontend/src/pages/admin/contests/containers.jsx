import { useParams } from 'react-router-dom';
import {
  getContestVictims,
  stopContestVictims,
  getContestTeams,
  getContestChallenges,
  downloadContainerTraffic,
  getContestVictimPods,
  getContestVictimPodLogs,
} from '../../../api/admin/contest';
import { getUserList } from '../../../api/admin/user';
import VictimInventory from '../../../components/features/Admin/victims/VictimInventory';
import ContestVictimStart from '../../../components/features/Admin/victims/ContestVictimStart';

export default function ContestContainers() {
  const { id } = useParams();
  const contestId = Number.parseInt(id, 10);
  const scope = {
    contestId,
    translationKey: 'admin.contests.containers',
    loadVictims: (params) => getContestVictims(contestId, params),
    stopVictims: (ids) => stopContestVictims(contestId, ids),
    loadPods: (victimId) => getContestVictimPods(contestId, victimId),
    loadLogs: (victimId, pod, container, lines) => getContestVictimPodLogs(contestId, victimId, pod, container, lines),
    downloadTraffic: (victim) => downloadContainerTraffic(contestId, victim.team_id, victim.id),
    search: {
      users: (name) => getUserList({ name, limit: 10, offset: 0 }),
      teams: (name) => getContestTeams(contestId, { name, limit: 10, offset: 0 }),
      challenges: (name) => getContestChallenges(contestId, { name, type: 'pods', limit: 10, offset: 0 }),
    },
  };

  return (
    <VictimInventory
      key={contestId}
      scope={scope}
      renderQuickActions={(refresh) => <ContestVictimStart contestId={contestId} onStarted={refresh} />}
    />
  );
}
