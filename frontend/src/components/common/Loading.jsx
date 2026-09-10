import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import Spinner from './Spinner';

function Loading({ compact = false, className = '' }) {
  const { t } = useTranslation();
  return (
    <div
      role="status"
      aria-live="polite"
      className={`flex items-center justify-center w-full ${compact ? 'min-h-32' : 'min-h-[500px]'} ${className}`}
    >
      <motion.div
        className="flex flex-col items-center gap-4"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.3 }}
      >
        <Spinner size="lg" border="md" />
        <span className="text-sm font-mono text-neutral-400">{t('common.loading')}</span>
      </motion.div>
    </div>
  );
}

export default Loading;
