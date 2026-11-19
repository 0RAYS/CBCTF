import { useParams } from 'react-router-dom';
import {
  getContestChallenges,
  getContestGeneratorLogs,
  getContestGenerators,
  getContestGeneratorStatus,
  startContestGenerators,
  stopContestGenerators,
} from '../../../api/admin/contest';
import GeneratorManagement from '../../../components/features/Admin/generators/GeneratorManagement';

export default function ContestGenerators() {
  const { id: contestId } = useParams();
  const api = {
    list: (params) => getContestGenerators(contestId, params),
    challenges: (params) => getContestChallenges(contestId, params),
    start: (challenges) => startContestGenerators(contestId, challenges),
    stop: (ids) => stopContestGenerators(contestId, ids),
    status: (id, signal) => getContestGeneratorStatus(contestId, id, signal),
    logs: (id, lines) => getContestGeneratorLogs(contestId, id, lines),
  };

  // key 随 contestId 变化重新挂载 GeneratorManagement
  return <GeneratorManagement key={contestId} api={api} />;
}
