import { useParams } from 'react-router-dom';
import CheatsManager from '../../../components/features/Admin/cheats/CheatsManager';

export default function AdminContestCheats() {
  const { id } = useParams();
  return <CheatsManager key={id} contestId={parseInt(id)} />;
}
