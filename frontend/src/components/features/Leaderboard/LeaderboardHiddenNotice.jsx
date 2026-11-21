import { motion } from 'motion/react';
import { IconEyeOff } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Card } from '@/components/common';

/** 排行榜黑灯提示组件 */
function LeaderboardHiddenNotice() {
  const { t } = useTranslation();

  return (
    <Card variant="default" padding="lg" animate>
      <div className="text-center space-y-4">
        <motion.div
          className="w-16 h-16 mx-auto border-2 border-yellow-400 rounded-full flex items-center justify-center"
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ delay: 0.2, type: 'spring' }}
        >
          <IconEyeOff size={32} className="text-yellow-400" />
        </motion.div>
        <h2 className="text-2xl font-mono text-neutral-50">{t('game.leaderboard.hidden.title')}</h2>
        <p className="text-neutral-400 max-w-xl mx-auto [overflow-wrap:anywhere]">
          {t('game.leaderboard.hidden.description')}
        </p>
      </div>
    </Card>
  );
}

export default LeaderboardHiddenNotice;
