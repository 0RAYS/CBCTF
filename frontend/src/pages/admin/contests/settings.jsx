import { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import Loading from '../../../components/common/Loading';

function AdminContestSettings() {
  const { id } = useParams();
  const navigate = useNavigate();

  useEffect(() => {
    // 重定向到主编辑页面
    navigate(`/admin/contests/${id}`);
  }, [id, navigate]);

  return <Loading />;
}

export default AdminContestSettings;
