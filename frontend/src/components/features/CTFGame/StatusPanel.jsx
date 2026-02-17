/**
 * 比赛状态面板组件
 * @param {Object} props
 * @param {Object} props.contestStatus - 比赛状态信息
 * @param {string} props.contestStatus.status - 比赛状态 ("upcoming" | "active" | "ended")
 * @param {string} props.contestStatus.startTime - 开始时间 (ISO格式)
 * @param {string} props.contestStatus.endTime - 结束时间 (ISO格式)
 * @param {Object} props.contestStatus.team - 当前队伍信息
 * @param {number} props.contestStatus.team.score - 队伍分数
 * @param {number} props.contestStatus.team.rank - 队伍排名
 * @param {number} props.contestStatus.team.solved - 已解题数
 * @param {number} props.contestStatus.totalChallenges - 总题目数
 * @example
 * const contestStatus = {
 *   status: "active",
 *   startTime: "2024-03-15T10:00:00Z",
 *   endTime: "2024-03-16T10:00:00Z",
 *   team: {
 *     score: 2450,
 *     rank: 12,
 *     solved: 12
 *   },
 *   totalChallenges: 30
 * }
 */

import { motion } from 'motion/react';
import { useState, useEffect } from 'react';
import InfoCard from './InfoCard';
import NotificationCard from './NotificationCard';
import { useTranslation } from 'react-i18next';

function StatusPanel({ contestStatus, notifications, onStatusExpired }) {
  const { t, i18n } = useTranslation();
  const [remainingSeconds, setRemainingSeconds] = useState(0);

  useEffect(() => {
    let intervalId;

    const updateRemainingTime = () => {
      const targetTime = Date.parse(contestStatus.endTime);
      const now = Date.now();
      const diff = targetTime - now;

      if (diff <= 0) {
        setRemainingSeconds(0);
        if (contestStatus.status === 'upcoming') {
          onStatusExpired?.('active');
        } else if (contestStatus.status === 'active') {
          onStatusExpired?.('ended');
        }
        return false;
      }

      setRemainingSeconds(Math.floor(diff / 1000));
      return true;
    };

    // 首次立即更新
    updateRemainingTime();

    // 使用setInterval每秒更新一次
    intervalId = setInterval(() => {
      updateRemainingTime();
    }, 1000);

    return () => {
      clearInterval(intervalId);
    };
  }, [contestStatus.status, contestStatus.startTime, contestStatus.endTime]);

  // 格式化显示
  const formatTime = (seconds) => {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
  };

  // 即将开始的状态
  if (contestStatus.status === 'upcoming') {
    return (
      <div className="grid grid-cols-4 gap-6">
        <motion.div
          className="col-span-3 grid grid-cols-3 gap-6"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <InfoCard
            title={t('game.statusPanel.contestStatus')}
            value={t('game.status.upcoming')}
            valueColor="text-yellow-400"
          />

          <InfoCard
            title={t('game.statusPanel.startTime')}
            value={new Date(contestStatus.startTime).toLocaleString(i18n.language || 'en-US')}
            valueColor="text-neutral-50"
          />

          <InfoCard
            title={t('game.statusPanel.timeUntilStart')}
            value={formatTime(remainingSeconds)}
            valueColor="text-neutral-50"
          />
        </motion.div>

        <NotificationCard notifications={notifications} />
      </div>
    );
  }
  if (contestStatus.status === 'ended') {
    return <></>;
  }

  // 比赛进行中的状态
  return (
    <div className="grid grid-cols-4 gap-6">
      <motion.div
        className="col-span-3 grid grid-cols-3 gap-6"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <InfoCard
          title={t('game.statusPanel.timeRemaining')}
          value={formatTime(remainingSeconds)}
          valueColor="text-neutral-50"
        />

        <InfoCard
          title={t('game.statusPanel.teamScore')}
          value={contestStatus.team.score.toLocaleString(i18n.language || 'en-US')}
          valueColor="text-geek-400"
          subTitle={t('common.rank')}
          subValue={`#${contestStatus.team.rank}`}
          subValueColor="text-yellow-400"
        />

        <InfoCard
          title={t('game.statusPanel.flagsSolved')}
          value={`${contestStatus.team.solved}/${contestStatus.totalChallenges}`}
          valueColor="text-neutral-50"
          showProgress
          progressValue={(contestStatus.team.solved / contestStatus.totalChallenges) * 100}
        />
      </motion.div>

      <NotificationCard notifications={notifications} />
    </div>
  );
}

export default StatusPanel;
