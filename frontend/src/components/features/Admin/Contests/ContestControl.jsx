import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconPlayerPause, IconPlayerPlay } from '@tabler/icons-react';
import { Button, Card } from '@/components/common';
import Checkbox from '@/components/common/Checkbox';
import { toast } from '@/utils/toast';
import { updateContestLeaderboard, updateContestPause } from '@/api/admin/contest';

export default function ContestControl({ contest, onChange }) {
  const { t } = useTranslation();
  const [pending, setPending] = useState(false);
  const [blackoutPending, setBlackoutPending] = useState(false);
  const paused = !!contest.paused;
  const leaderboardHidden = !!contest.leaderboardHidden;

  const togglePause = async () => {
    if (pending) return;
    setPending(true);
    try {
      const response = await updateContestPause(contest.id, !paused);
      if (response.code === 200) {
        toast.success({
          description: t(paused ? 'admin.contests.control.toast.resumed' : 'admin.contests.control.toast.paused'),
        });
        onChange({ paused: !paused });
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.control.toast.failed') });
    } finally {
      setPending(false);
    }
  };

  const toggleBlackout = async (event) => {
    const hidden = event.target.checked;
    if (blackoutPending) return;
    setBlackoutPending(true);
    try {
      const response = await updateContestLeaderboard(contest.id, hidden);
      if (response.code === 200) {
        toast.success({
          description: t(
            hidden ? 'admin.contests.control.toast.blackoutEnabled' : 'admin.contests.control.toast.blackoutDisabled'
          ),
        });
        onChange({ leaderboardHidden: hidden });
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.control.toast.failed') });
    } finally {
      setBlackoutPending(false);
    }
  };

  return (
    <Card padding="lg" className="space-y-3">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div className="min-w-0">
          <h3 className="text-lg font-mono text-neutral-50">{t('admin.contests.control.title')}</h3>
          <p className="mt-1 text-sm text-neutral-500">{t('admin.contests.control.help')}</p>
        </div>
        <Button
          variant={paused ? 'primary' : 'secondary'}
          size="sm"
          disabled={pending}
          onClick={togglePause}
          icon={paused ? <IconPlayerPlay size={16} /> : <IconPlayerPause size={16} />}
        >
          {t(paused ? 'admin.contests.control.resume' : 'admin.contests.control.pause')}
        </Button>
      </div>
      <p className="text-sm text-neutral-400">
        {t(paused ? 'admin.contests.control.pausedHint' : 'admin.contests.control.runningHint')}
      </p>
      <div className="flex flex-wrap items-center justify-between gap-4 border-t border-neutral-300/20 pt-3">
        <div className="min-w-0">
          <p className="text-sm text-neutral-50">{t('admin.contests.control.blackout.label')}</p>
          <p className="mt-1 text-sm text-neutral-500">{t('admin.contests.control.blackout.help')}</p>
        </div>
        <Checkbox
          checked={leaderboardHidden}
          onChange={toggleBlackout}
          disabled={blackoutPending}
          label={t('admin.contests.control.blackout.toggle')}
          className="text-sm text-neutral-300"
        />
      </div>
    </Card>
  );
}
