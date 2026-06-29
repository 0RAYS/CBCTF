import { useState } from 'react';
import { useLocation, Link } from 'react-router-dom';
import { motion, AnimatePresence } from 'motion/react';
import { useTranslation } from 'react-i18next';
import { Button } from '../../components/common';
import { verifyEmail } from '../../api/auth';

function VerifyEmail() {
  const { t } = useTranslation();
  const { search } = useLocation();

  const params = new URLSearchParams(search);
  const token = params.get('token') ?? '';
  const paramsValid = Boolean(token);

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [succeeded, setSucceeded] = useState(false);
  const [tokenInvalid, setTokenInvalid] = useState(!paramsValid);

  const handleActivate = async () => {
    setIsSubmitting(true);
    try {
      const response = await verifyEmail({ token });
      if (response.code === 200) {
        setSucceeded(true);
      } else {
        setTokenInvalid(true);
      }
    } catch {
      setTokenInvalid(true);
    } finally {
      setIsSubmitting(false);
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
        {/* 标题 */}
        <div className="relative flex justify-center mb-8">
          <div className="absolute -top-[3px] -right-[3px] w-[10px] h-[10px] border-t border-r border-neutral-300" />
          <div className="absolute -bottom-[3px] -left-[3px] w-[10px] h-[10px] border-b border-l border-neutral-300" />
          <h1 className="text-neutral-50 text-2xl font-mono tracking-wider">{t('auth.verifyEmail.title')}</h1>
        </div>

        <AnimatePresence mode="wait">
          {/* token 无效 */}
          {tokenInvalid && (
            <motion.div
              key="invalid"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="space-y-4 text-center"
            >
              <p className="text-red-400 text-sm font-mono leading-relaxed">{t('auth.verifyEmail.invalidToken')}</p>
              <Link to="/settings">
                <Button variant="outline" fullWidth>
                  {t('auth.verifyEmail.requestNew')}
                </Button>
              </Link>
            </motion.div>
          )}

          {/* 激活成功 */}
          {!tokenInvalid && succeeded && (
            <motion.div
              key="success"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="space-y-4 text-center"
            >
              <p className="text-neutral-300 text-sm leading-relaxed">{t('auth.verifyEmail.successMessage')}</p>
              <Link to="/settings">
                <Button variant="primary" fullWidth>
                  {t('auth.verifyEmail.goToSettings')}
                </Button>
              </Link>
            </motion.div>
          )}

          {/* 确认激活 */}
          {!tokenInvalid && !succeeded && (
            <motion.div
              key="confirm"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="space-y-6"
            >
              {/* 账号确认区域 */}
              <div className="rounded-md border border-neutral-600/60 bg-neutral-900/50 px-5 py-4 space-y-2">
                <p className="text-neutral-400 text-xs font-mono uppercase tracking-widest">
                  {t('auth.verifyEmail.title')}
                </p>
                <p className="text-neutral-100 text-base font-mono">{t('auth.verifyEmail.question')}</p>
                <p className="text-neutral-500 text-xs leading-relaxed">{t('auth.verifyEmail.description')}</p>
              </div>

              <Button variant="primary" fullWidth disabled={isSubmitting} onClick={handleActivate}>
                {isSubmitting ? t('auth.verifyEmail.submitting') : t('auth.verifyEmail.submit')}
              </Button>
            </motion.div>
          )}
        </AnimatePresence>
      </motion.div>
    </div>
  );
}

export default VerifyEmail;
