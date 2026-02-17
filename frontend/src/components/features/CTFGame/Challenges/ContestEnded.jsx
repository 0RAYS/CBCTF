import { motion } from 'motion/react';
import { Button, Card } from '../../../common';
import { useTranslation } from 'react-i18next';
import WriteupUpload from './WriteupUpload';

function ContestEnded({ contestInfo, onViewScoreboard, onUploadWriteup, onViewChallenges, writeups = [] }) {
  const { t } = useTranslation();

  return (
    <div className="w-full space-y-6">
      <Card variant="default" padding="lg" animate>
        {/* 标题和图标 */}
        <div className="text-center space-y-4 mb-8">
          <motion.div
            className="w-16 h-16 mx-auto border-2 border-yellow-400 rounded-full
                            flex items-center justify-center"
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: 0.2, type: 'spring' }}
          >
            <span className="text-3xl">🏆</span>
          </motion.div>
          <h2 className="text-2xl font-mono text-neutral-50">
            {t('game.contestEnded.titlePrefix')}{' '}
            <span className="text-yellow-400">{t('game.contestEnded.titleHighlight')}</span>
          </h2>
        </div>

        {/* 比赛信息 */}
        <div className="space-y-6">
          {/* 比赛统计 */}
          <div className="grid grid-cols-3 gap-4">
            <motion.div
              className="p-4 border border-neutral-300/30 rounded-md text-center"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 }}
            >
              <div className="text-neutral-400 text-sm mb-1">{t('common.duration')}</div>
              <div className="text-neutral-50 font-mono">{contestInfo?.duration || '24h'}</div>
            </motion.div>
            <motion.div
              className="p-4 border border-neutral-300/30 rounded-md text-center"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.4 }}
            >
              <div className="text-neutral-400 text-sm mb-1">{t('common.teams')}</div>
              <div className="text-neutral-50 font-mono">{contestInfo?.totalTeams || '256'}</div>
            </motion.div>
            <motion.div
              className="p-4 border border-neutral-300/30 rounded-md text-center"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.5 }}
            >
              <div className="text-neutral-400 text-sm mb-1">{t('common.challenges')}</div>
              <div className="text-neutral-50 font-mono">{contestInfo?.totalChallenges || '30'}</div>
            </motion.div>
          </div>

          {/* 团队成绩 */}
          <motion.div
            className="p-6 border border-geek-400 rounded-md bg-geek-400/5"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.6 }}
          >
            <div className="flex items-center justify-between mb-4">
              <span className="text-neutral-50 font-mono">{t('game.contestEnded.teamPerformance.title')}</span>
              <span className="text-geek-400 font-mono">#{contestInfo?.teamRank || '1'}</span>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <div className="text-neutral-400 text-sm mb-1">{t('common.score')}</div>
                <div className="text-geek-400 font-mono text-xl">{contestInfo?.teamScore || '0'}</div>
              </div>
              <div>
                <div className="text-neutral-400 text-sm mb-1">{t('common.solved')}</div>
                <div className="text-neutral-50 font-mono text-xl">
                  {contestInfo?.teamSolved || '0'}/{contestInfo?.totalChallenges || '0'}
                </div>
              </div>
            </div>
          </motion.div>

          {/* 操作按钮 */}
          <motion.div
            className="flex justify-center gap-4 pt-4"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.7 }}
          >
            <Button variant="outline" size="lg" onClick={onViewChallenges}>
              {t('game.contestEnded.actions.viewChallenges')}
            </Button>
            <Button variant="primary" size="lg" onClick={onViewScoreboard}>
              {t('game.contestEnded.actions.viewFinalScoreboard')}
            </Button>
          </motion.div>
        </div>
      </Card>

      {/* 题解上传区域 */}
      <WriteupUpload onUploadWriteup={onUploadWriteup} writeups={writeups} />
    </div>
  );
}

export default ContestEnded;
