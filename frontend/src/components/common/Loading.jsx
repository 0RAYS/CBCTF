import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';

function Loading() {
  const { t } = useTranslation();
  return (
    <div className="flex items-center justify-center min-h-[500px] w-full">
      <motion.div
        className="text-neutral-300 font-mono text-lg"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.5 }}
      >
        {t('common.loading')}
      </motion.div>
    </div>
    // <div className="min-h-[400px] flex items-center justify-center">
    //   <div className="flex flex-col items-center gap-4">
    //     <div className="w-8 h-8 border-2 border-t-geek-400 border-neutral-400/30 rounded-full animate-spin" />
    //     <span className="text-sm font-mono text-neutral-400">加载中...</span>
    //   </div>
    // </div>
  );
}

export default Loading;
