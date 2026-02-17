import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { fetchUserInfo } from '../store/user';
import { toast } from '../utils/toast';
import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';

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

      // 从URL参数中获取token
      const token = searchParams.get('token');
      const error = searchParams.get('error');

      if (error) {
        throw new Error(t('oauth.callback.errorWithReason', { reason: error }));
      }

      if (!token) {
        throw new Error(t('oauth.callback.tokenMissing'));
      }

      // 将token保存到localStorage
      localStorage.setItem('token', 'Bearer ' + token);
      localStorage.setItem('userType', 'user');

      // 设置请求头中的Authorization
      // 这里需要更新axios的默认请求头
      // 由于request.js中可能已经有token处理逻辑，我们只需要确保localStorage中有token即可

      // 等待一下确保token已保存
      await new Promise((resolve) => setTimeout(resolve, 500));

      // 获取用户信息
      await dispatch(fetchUserInfo());

      setStatus('success');
      toast.success({ description: t('oauth.callback.toast.success') });

      // 跳转到游戏页面
      navigate('/games');
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
