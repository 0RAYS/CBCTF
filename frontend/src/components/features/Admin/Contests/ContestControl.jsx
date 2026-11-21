import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { IconPlayerPause, IconPlayerPlay } from '@tabler/icons-react';
import { Button, Card } from '@/components/common';
import { toast } from '@/utils/toast';
import { updateContestPause } from '@/api/admin/contest';

export default function ContestControl({ contest, onChange }) {
  const { t } = useTranslation();
  const [pending, setPending] = useState(false);
  const paused = !!contest.paused;

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
    </Card>
  );
}
