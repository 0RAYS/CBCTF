import {
  getVictims,
  stopVictims,
  getVictimTraffic,
  downloadVictimTraffic,
  getVictimPods,
  getVictimPodLogs,
} from '../../api/admin/victims';
import { getUserList } from '../../api/admin/user';
import { getChallengeList } from '../../api/admin/challenge';
import VictimInventory from '../../components/features/Admin/victims/VictimInventory';

const scope = {
  translationKey: 'admin.victims',
  loadVictims: getVictims,
  stopVictims,
  loadPods: getVictimPods,
  loadLogs: getVictimPodLogs,
  fetchTraffic: (victim, params) => getVictimTraffic(victim.id, params),
  downloadTraffic: (victim) => downloadVictimTraffic(victim.id),
  search: {
    users: (name) => getUserList({ name, limit: 10, offset: 0 }),
    challenges: (name) => getChallengeList({ name, type: 'pods', limit: 10, offset: 0 }),
  },
};

export default function AdminVictims() {
  return <VictimInventory scope={scope} />;
}
