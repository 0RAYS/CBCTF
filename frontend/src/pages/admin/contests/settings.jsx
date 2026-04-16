import { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

function AdminContestSettings() {
  const { id } = useParams();
  const navigate = useNavigate();

  useEffect(() => {
    // 重定向到主编辑页面, 两个页面已合并
    navigate(`/admin/contests/${id}`);
  }, [id, navigate]);

  return null;
}

export default AdminContestSettings;
