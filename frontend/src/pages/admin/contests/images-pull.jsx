import { useParams } from 'react-router-dom';
import { detectContestImages, getContestPullImages, pullContestImages } from '@/api/admin/contest';
import ImagesPullManagement from '@/components/features/Admin/images/ImagesPullManagement.jsx';

function AdminContestImagesPull() {
  const { id } = useParams();

  return (
    <ImagesPullManagement
      key={id}
      scope="contest"
      fetchImages={() => getContestPullImages(parseInt(id, 10))}
      pullImages={(data) => pullContestImages(parseInt(id, 10), data)}
      detectImages={() => detectContestImages(parseInt(id, 10))}
    />
  );
}

export default AdminContestImagesPull;
