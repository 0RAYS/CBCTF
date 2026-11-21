import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { fetchAccessibleRoutes, fetchUserInfo } from '@/store/user';
import { store } from '@/store';
import { toast } from '@/utils/toast';
import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import { IconCheck, IconX } from '@tabler/icons-react';
import { Spinner } from '@/components/common';
import { exchangeOauthCode } from '@/api/oauth';

function StatusHeading({ children }) {
  return (
    <div className="relative flex justify-center w-full">
      <div className="absolute -top-[3px] -right-[3px] w-[10px] h-[10px] border-t border-r border-neutral-300" />
      <div className="absolute -bottom-[3px] -left-[3px] w-[10px] h-[10px] border-b border-l border-neutral-300" />
      <h1 className="text-2xl font-mono text-neutral-50 tracking-wider">{children}</h1>
    </div>
  );
}

function OAuthCallback() {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState('processing'); // 'processing' | 'success' | 'error'

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

      // 用一次性 code 换取真实 token（服务端写入 httpOnly cookie）
      const res = await exchangeOauthCode(code);
      if (res?.code !== 200) {
        throw new Error(t('oauth.callback.tokenMissing'));
      }

      // 获取用户信息及可访问路由（需两者都完成才能正确判断权限）
      await dispatch(fetchUserInfo());
      await dispatch(fetchAccessibleRoutes());

      setStatus('success');
      toast.success({ description: t('oauth.callback.toast.success') });

      // 跳转规则与 Login.jsx 保持一致
      const { hasAdminAccess } = store.getState().user;
      navigate(hasAdminAccess ? '/admin/dashboard' : '/games');
    } catch (error) {
      setStatus('error');
      toast.danger({ description: error.message || t('oauth.callback.toast.failed') });

      setTimeout(() => {
        navigate('/login');
      }, 2000);
    }
  };

  useEffect(() => {
    handleOAuthCallback();
  }, []);

  const renderContent = () => {
    switch (status) {
      case 'processing':
        return (
          <div className="flex flex-col items-center justify-center space-y-4">
            <Spinner size="lg" border="md" />
            <StatusHeading>{t('oauth.callback.processing.title')}</StatusHeading>
            <p className="text-neutral-400 text-sm">{t('oauth.callback.processing.subtitle')}</p>
          </div>
        );

      case 'success':
        return (
          <div className="flex flex-col items-center justify-center space-y-4">
            <div className="w-12 h-12 bg-green-400/20 text-green-400 rounded-full flex items-center justify-center">
              <IconCheck size={24} />
            </div>
            <StatusHeading>{t('oauth.callback.success.title')}</StatusHeading>
            <p className="text-neutral-400 text-sm">{t('oauth.callback.success.subtitle')}</p>
          </div>
        );

      case 'error':
        return (
          <div className="flex flex-col items-center justify-center space-y-4">
            <div className="w-12 h-12 bg-red-400/20 text-red-400 rounded-full flex items-center justify-center">
              <IconX size={24} />
            </div>
            <StatusHeading>{t('oauth.callback.error.title')}</StatusHeading>
            <p className="text-neutral-400 text-sm">{t('oauth.callback.error.subtitle')}</p>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className="flex items-center justify-center min-h-[calc(100vh-190px)] px-4">
      <motion.div
        className="w-full max-w-[400px] bg-neutral-800/80 border border-neutral-600/60 rounded-md p-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3, ease: [0.25, 1, 0.5, 1] }}
      >
        {renderContent()}
      </motion.div>
    </div>
  );
}

export default OAuthCallback;
