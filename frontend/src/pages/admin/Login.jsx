import { useNavigate } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { adminLogin } from '../../api/auth';
import { fetchUserInfo } from '../../store/user';
import { setEffectLoading } from '../../store/app';
import AdminAuthPanel from '../../components/features/Auth/AdminAuthPanel';
import { toast } from '../../utils/toast';
import { useTranslation } from 'react-i18next';

function AdminLogin() {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { t } = useTranslation();

  const handleAuth = async ({ data }) => {
    try {
      dispatch(setEffectLoading({ effect: 'adminLogin', loading: true }));
      const response = await adminLogin({
        name: data.username,
        password: data.password,
      });

      if (response.code === 200) {
        await dispatch(fetchUserInfo(true));
        toast.success({ description: t('toast.auth.loginSuccess') });
        navigate('/admin/dashboard');
      }
    } catch (error) {
      toast.danger({ title: t('toast.auth.loginFailed'), description: error.message || t('toast.auth.loginFailed') });
    } finally {
      dispatch(setEffectLoading({ effect: 'adminLogin', loading: false }));
    }
  };

  return (
    <div className="flex items-center justify-center min-h-[calc(100vh-190px)]">
      <AdminAuthPanel onSubmit={handleAuth} />
    </div>
  );
}

export default AdminLogin;
