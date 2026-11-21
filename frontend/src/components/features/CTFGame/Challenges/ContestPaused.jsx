import { motion } from 'motion/react';
import { IconPlayerPause } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Button, Card } from '@/components/common';

function ContestPaused({ onViewLeaderboard }) {
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
          <IconPlayerPause size={32} className="text-yellow-400" />
        </motion.div>
        <h2 className="text-2xl font-mono text-neutral-50">{t('game.contestPaused.title')}</h2>
        <p className="text-neutral-400 max-w-xl mx-auto [overflow-wrap:anywhere]">
          {t('game.contestPaused.description')}
        </p>
      </div>
      <div className="flex justify-center gap-4 pt-6">
        <Button variant="primary" size="lg" onClick={onViewLeaderboard}>
          {t('game.contestPaused.actions.viewLeaderboard')}
        </Button>
      </div>
    </Card>
  );
}

export default ContestPaused;
