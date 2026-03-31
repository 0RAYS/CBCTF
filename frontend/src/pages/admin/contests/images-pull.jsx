import { useParams } from 'react-router-dom';
import { getContestPullImages, pullContestImages } from '../../../api/admin/contest';
import AdminImagesPullPage from '../../../components/features/Admin/AdminImagesPullPage.jsx';

function AdminContestImagesPull() {
  const { id } = useParams();

  return (
    <AdminImagesPullPage
      scope="contest"
      refreshKey={id}
      fetchImages={() => getContestPullImages(parseInt(id, 10))}
      pullImages={(data) => pullContestImages(parseInt(id, 10), data)}
    />
  );
}

export default AdminContestImagesPull;
