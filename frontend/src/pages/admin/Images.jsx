import { getAdminPullImages, pullAdminImages } from '../../api/admin/image';
import ImagesPullManagement from '../../components/features/Admin/images/ImagesPullManagement.jsx';

function AdminImages() {
  return <ImagesPullManagement scope="global" fetchImages={getAdminPullImages} pullImages={pullAdminImages} />;
}

export default AdminImages;
