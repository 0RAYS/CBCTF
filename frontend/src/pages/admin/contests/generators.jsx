import { useParams } from 'react-router-dom';
import {
  getContestGenerators,
  startContestGenerators,
  stopContestGenerators,
  getContestChallenges,
  getContestGeneratorLogs,
} from '../../../api/admin/contest';
import GeneratorManagement from '../../../components/features/Admin/generators/GeneratorManagement';

export default function ContestGenerators() {
  const { id: contestId } = useParams();
  const api = {
    list: (params) => getContestGenerators(contestId, params),
    challenges: (params) => getContestChallenges(contestId, params),
    start: (challenges) => startContestGenerators(contestId, challenges),
    stop: (ids) => stopContestGenerators(contestId, ids),
    logs: (id, lines) => getContestGeneratorLogs(contestId, id, lines),
  };

  // A contest change owns a new query, operation, and dialog session.
  return <GeneratorManagement key={contestId} api={api} textKey="admin.contests.generators" />;
}
