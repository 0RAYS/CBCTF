import { motion } from 'motion/react';
import { useEffect, useState } from 'react';
import InfoCard from '@/components/features/CTFGame/InfoCard';
import NotificationCard from '@/components/features/CTFGame/NotificationCard';
import { useTranslation } from 'react-i18next';

/** 比赛状态面板组件 */
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

    updateRemainingTime();

    intervalId = setInterval(() => {
      updateRemainingTime();
    }, 1000);

    return () => {
      clearInterval(intervalId);
    };
  }, [contestStatus.status, contestStatus.startTime, contestStatus.endTime]);

  const formatTime = (seconds) => {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
  };

  if (contestStatus.status === 'upcoming') {
    return (
      <div className="grid grid-cols-1 xl:grid-cols-4 gap-4 [overflow-wrap:anywhere] [&>div]:min-w-0">
        <motion.div
          className="xl:col-span-3 grid grid-cols-1 sm:grid-cols-3 gap-4 [&>div]:min-w-0 [&>div>div.flex]:flex-col [&>div>div.flex]:items-start [&>div>div.flex]:gap-2"
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

  return (
    <div className="grid grid-cols-1 xl:grid-cols-4 gap-4 [overflow-wrap:anywhere] [&>div]:min-w-0">
      <motion.div
        className="xl:col-span-3 grid grid-cols-1 sm:grid-cols-3 gap-4 [&>div]:min-w-0 [&>div>div.flex]:flex-col [&>div>div.flex]:items-start [&>div>div.flex]:gap-2"
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
