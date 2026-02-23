import { useNavigate } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { login, register } from '../api/auth';
import { fetchUserInfo, fetchAccessibleRoutes } from '../store/user';
import { setEffectLoading } from '../store/app';
import AuthPanel from '../components/features/Auth/AuthPanel';
import { toast } from '../utils/toast.js';
import { useTranslation } from 'react-i18next';

function Login() {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { t } = useTranslation();

  const handleAuth = async ({ type, data }) => {
    const msg = type === 'login' ? t('toast.auth.loginFailed') : t('toast.auth.registerFailed');
    try {
      let response;
      dispatch(setEffectLoading({ effect: type, loading: true }));
      if (type === 'login') {
        response = await login({
          name: data.username,
          password: data.password,
        });
      } else {
        // 处理注册逻辑
        response = await register({
          name: data.username,
          email: data.email,
          password: data.password,
        });
      }
      if (response.code === 200) {
        await dispatch(fetchUserInfo());
        await dispatch(fetchAccessibleRoutes());
        // 根据角色跳转到不同首页
        const state = await import('../store/index').then((m) => m.default.getState());
        navigate(state.user.isAdmin ? '/admin/dashboard' : '/games');
      }
    } catch (error) {
      toast.danger({ title: msg, description: error.message });
      throw error;
    } finally {
      dispatch(setEffectLoading({ effect: type, loading: false }));
    }
  };

  return (
    <div className="flex items-center justify-center min-h-[calc(100vh-190px)]">
      <AuthPanel onSubmit={handleAuth} />
    </div>
  );
}

export default Login;
