import {
  getGeneratorLogs,
  getGenerators,
  getGeneratorStatus,
  startGenerators,
  stopGenerators,
} from '../../api/admin/generators';
import { getChallengeList } from '../../api/admin/challenge';
import GeneratorManagement from '../../components/features/Admin/generators/GeneratorManagement';

const api = {
  list: getGenerators,
  challenges: getChallengeList,
  start: startGenerators,
  stop: stopGenerators,
  logs: getGeneratorLogs,
  status: getGeneratorStatus,
};

export default function AdminGenerators() {
  return <GeneratorManagement api={api} textKey="admin.generators" />;
}
