import { useState, useEffect } from 'react';
import { motion } from 'motion/react';
import { getOAuthProviders } from '../../../api/oauth';
import { toast } from '../../../utils/toast';
import { useTranslation } from 'react-i18next';

function OAuthLogin() {
  const [providers, setProviders] = useState({});
  const [loading, setLoading] = useState(true);
  const { t } = useTranslation();

  useEffect(() => {
    fetchOAuthProviders();
  }, []);

  const fetchOAuthProviders = async () => {
    try {
      setLoading(true);
      const response = await getOAuthProviders();
      if (response.code === 200) {
        setProviders(response.data);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('toast.oauth.providersFailed') });
    } finally {
      setLoading(false);
    }
  };

  const handleOAuthLogin = (providerName, loginUrl) => {
    // 跳转到OAuth登录页面
    window.location.href = loginUrl;
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center py-4">
        <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-neutral-300"></div>
      </div>
    );
  }

  // 如果没有可用的OAuth提供商，不显示任何内容
  if (Object.keys(providers).length === 0) {
    return null;
  }

  return (
    <motion.div
      className="space-y-3"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, delay: 0.2 }}
    >
      {/* 分割线 */}
      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-neutral-600"></div>
        </div>
        <div className="relative flex justify-center text-xs">
          <span className="bg-neutral-900 px-2 text-neutral-400">{t('auth.oauth.or')}</span>
        </div>
      </div>

      {/* OAuth按钮 */}
      <div className="space-y-2">
        {Object.entries(providers).map(([providerName, provider]) => (
          <motion.button
            key={providerName}
            onClick={() => handleOAuthLogin(providerName, provider.url)}
            className="w-full h-[40px] bg-black/20 border border-neutral-300 rounded-md px-4
                     text-neutral-50 hover:border-geek-400 hover:shadow-focus
                     transition-all duration-200 flex items-center justify-center gap-3"
            whileHover={{ opacity: 0.85 }}
            whileTap={{ opacity: 0.7 }}
          >
            <img
              src={provider.picture}
              alt={provider.name}
              loading="lazy"
              className="w-5 h-5 rounded-full"
              onError={(e) => {
                e.target.style.display = 'none';
              }}
            />
            <span className="font-medium">{t('auth.oauth.useProvider', { provider: provider.name })}</span>
          </motion.button>
        ))}
      </div>
    </motion.div>
  );
}

export default OAuthLogin;
