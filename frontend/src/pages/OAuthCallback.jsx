import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { fetchUserInfo, fetchAccessibleRoutes } from '../store/user';
import { store } from '../store';
import { toast } from '../utils/toast';
import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import { exchangeOauthCode } from '../api/oauth';

function OAuthCallback() {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState('processing'); // 'processing' | 'success' | 'error'

  useEffect(() => {
    handleOAuthCallback();
  }, []);

  const handleOAuthCallback = async () => {
    try {
      setStatus('processing');

      const code = searchParams.get('code');
      const error = searchParams.get('error');

      if (error) {
        throw new Error(t('oauth.callback.errorWithReason', { reason: error }));
      }

      if (!code) {
        throw new Error(t('oauth.callback.tokenMissing'));
      }

      // 用一次性 code 换取真实 token
      const res = await exchangeOauthCode(code);
      const token = res?.data?.token;
      if (!token) {
        throw new Error(t('oauth.callback.tokenMissing'));
      }

      // 将token保存到localStorage
      localStorage.setItem('token', 'Bearer ' + token);
      localStorage.setItem('userType', 'user');

      // 等待一下确保token已保存
      await new Promise((resolve) => setTimeout(resolve, 500));

      // 获取用户信息及可访问路由（需两者都完成才能正确判断权限）
      await dispatch(fetchUserInfo());
      await dispatch(fetchAccessibleRoutes());

      setStatus('success');
      toast.success({ description: t('oauth.callback.toast.success') });

      // 根据权限跳转（与 Login.jsx 保持一致）
      const { hasAdminAccess } = store.getState().user;
      navigate(hasAdminAccess ? '/admin/dashboard' : '/games');
    } catch (error) {
      setStatus('error');
      toast.danger({ description: error.message || t('oauth.callback.toast.failed') });

      // 延迟跳转到登录页面
      setTimeout(() => {
        navigate('/login');
      }, 2000);
    }
  };

  const renderContent = () => {
    switch (status) {
      case 'processing':
        return (
          <div className="flex flex-col items-center justify-center space-y-4">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-geek-400"></div>
            <p className="text-neutral-300 text-lg">{t('oauth.callback.processing.title')}</p>
            <p className="text-neutral-400 text-sm">{t('oauth.callback.processing.subtitle')}</p>
          </div>
        );

      case 'success':
        return (
          <div className="flex flex-col items-center justify-center space-y-4">
            <div className="w-12 h-12 bg-green-500 rounded-full flex items-center justify-center">
              <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <p className="text-neutral-300 text-lg">{t('oauth.callback.success.title')}</p>
            <p className="text-neutral-400 text-sm">{t('oauth.callback.success.subtitle')}</p>
          </div>
        );

      case 'error':
        return (
          <div className="flex flex-col items-center justify-center space-y-4">
            <div className="w-12 h-12 bg-red-500 rounded-full flex items-center justify-center">
              <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </div>
            <p className="text-neutral-300 text-lg">{t('oauth.callback.error.title')}</p>
            <p className="text-neutral-400 text-sm">{t('oauth.callback.error.subtitle')}</p>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="flex items-center justify-center min-h-[calc(100vh-190px)]">
      <motion.div
        className="w-[400px] bg-black/40 backdrop-blur-[2px] border border-neutral-300 rounded-md p-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
      >
        {renderContent()}
      </motion.div>
    </div>
  );
}

export default OAuthCallback;
