/**
 * 积分榜总览统计组件
 * @param {Object} props
 * @param {number} props.totalTeams - 总队伍数
 * @param {number} props.totalSolves - 总解题数
 * @param {number} props.highestScore - 最高分数
 * @param {number} props.totalPlayers - 总人数
 */

import { motion } from 'motion/react';
import { Card } from '../../common';
import { useTranslation } from 'react-i18next';

function ScoreboardStats({ totalTeams, totalSolves, highestScore, totalPlayers }) {
  const { t, i18n } = useTranslation();

  const stats = [
    { label: t('game.scoreboardStats.totalTeams'), value: totalTeams, valueClass: 'text-neutral-50' },
    { label: t('game.scoreboardStats.totalPlayers'), value: totalPlayers, valueClass: 'text-neutral-50' },
    { label: t('game.scoreboardStats.totalSolves'), value: totalSolves, valueClass: 'text-geek-400' },
    { label: t('game.scoreboardStats.highestScore'), value: highestScore, valueClass: 'text-yellow-400' },
  ];

  return (
    <div className="grid grid-cols-2 gap-3 lg:grid-cols-4 lg:gap-6">
      {stats.map((stat, index) => (
        <Card key={index} variant="default" padding="md" className="min-w-0">
          <div className="text-neutral-400 text-xs sm:text-sm leading-relaxed">{stat.label}</div>
          <motion.div
            className={`mt-1 font-mono tabular-nums text-xl sm:text-2xl [overflow-wrap:anywhere] ${stat.valueClass}`}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.2, delay: index * 0.05 }}
          >
            {typeof stat.value === 'number' ? stat.value.toLocaleString(i18n.language || 'en-US') : stat.value}
          </motion.div>
        </Card>
      ))}
    </div>
  );
}

export default ScoreboardStats;
