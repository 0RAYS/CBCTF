import { getAdminPullImages, pullAdminImages } from '../../api/admin/image';
import AdminImagesPullPage from '../../components/features/Admin/AdminImagesPullPage.jsx';

function AdminImages() {
  return <AdminImagesPullPage scope="global" fetchImages={getAdminPullImages} pullImages={pullAdminImages} />;
}

export default AdminImages;
