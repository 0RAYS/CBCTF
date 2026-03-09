import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';

function Loading() {
  const { t } = useTranslation();
  return (
    <div role="status" aria-live="polite" className="flex items-center justify-center min-h-[500px] w-full">
      <motion.div
        className="flex flex-col items-center gap-4"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.3 }}
      >
        <div className="w-8 h-8 border-2 border-t-geek-400 border-neutral-400/30 rounded-full animate-spin" />
        <span className="text-sm font-mono text-neutral-400">{t('common.loading')}</span>
      </motion.div>
    </div>
  );
}

export default Loading;
