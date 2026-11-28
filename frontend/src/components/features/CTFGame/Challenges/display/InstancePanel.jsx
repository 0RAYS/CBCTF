import { useTranslation } from 'react-i18next';
import { IconCheck, IconCopy, IconDownload } from '@tabler/icons-react';
import { Button } from '@/components/common';
import {
  formatTimeLeft,
  normalizeInstanceStatus,
  parseSshEndpoint,
} from '@/components/features/CTFGame/Challenges/models/challengeViewModel';

export default function InstancePanel({
  challenge,
  timeLeft,
  loading,
  onLaunch,
  onExtend,
  onDestroy,
  onCopy,
  isCopied,
}) {
  const { t } = useTranslation();
  const status = normalizeInstanceStatus(challenge.instanceStatus);
  const isRunning = status === 'running';
  const isWaiting = status === 'waiting';
  const isPending = status === 'pending';
  const isTerminating = status === 'terminating';
  const duration = Number(challenge.instanceDuration) || 0;
  const progressWidth = duration > 0 ? Math.max(0, Math.min(100, (timeLeft / duration) * 100)) : 0;
  const accessKey = challenge.instanceAccessKey || '';
  const accessUser = challenge.instanceAccessUser || '';
  const sshEndpoint = parseSshEndpoint(challenge.instanceIP || []);
  const sshCommand = sshEndpoint ? `ssh -i id_ed25519 -p ${sshEndpoint.port} ${accessUser}@${sshEndpoint.host}` : '';
  const launchButtonLabel = isWaiting
    ? t('game.challengeModal.instance.waiting')
    : isTerminating
      ? t('game.challengeModal.instance.terminating')
      : isPending
        ? t('game.challengeModal.actions.launching')
        : t('game.challengeModal.actions.launch');

  const downloadKey = () => {
    const url = URL.createObjectURL(new Blob([accessKey], { type: 'application/octet-stream' }));
    const link = document.createElement('a');
    link.href = url;
    link.download = 'id_ed25519';
    link.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="space-y-3">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex flex-wrap items-center gap-4">
          <div className="flex items-center gap-2">
            <span
              className={`w-2 h-2 rounded-full transition-colors duration-300 ${
                isRunning
                  ? 'bg-green-400'
                  : isTerminating
                    ? 'bg-orange-400'
                    : isPending || isWaiting
                      ? 'bg-yellow-400'
                      : 'bg-neutral-500'
              }`}
            />
            <span className="text-neutral-50 font-mono text-sm">
              {t(`game.challengeModal.instance.${status || 'notRunning'}`)}
            </span>
          </div>
          {isRunning && (
            <div className="flex items-center gap-1.5">
              <span className="text-neutral-400 text-xs">{t('game.challengeModal.instance.time')}</span>
              <span className="text-yellow-400 font-mono text-sm">{formatTimeLeft(timeLeft)}</span>
            </div>
          )}
        </div>
        <div className="flex flex-wrap items-center gap-2">
          {!isRunning ? (
            <Button
              variant="primary"
              size="sm"
              onClick={onLaunch}
              disabled={loading.launching || isWaiting || isPending || isTerminating}
              loading={loading.launching}
              className={isWaiting || isPending || isTerminating ? 'border-yellow-400 text-yellow-400' : ''}
            >
              {launchButtonLabel}
            </Button>
          ) : (
            <>
              <Button
                variant="primary"
                size="sm"
                onClick={onExtend}
                disabled={loading.extending}
                loading={loading.extending}
                className="border-yellow-400 text-yellow-400 hover:bg-yellow-400/10"
              >
                {loading.extending
                  ? t('game.challengeModal.actions.extending')
                  : t('game.challengeModal.actions.extend')}
              </Button>
              <Button
                variant="danger"
                size="sm"
                onClick={onDestroy}
                disabled={loading.destroying}
                loading={loading.destroying}
              >
                {loading.destroying
                  ? t('game.challengeModal.actions.destroying')
                  : t('game.challengeModal.actions.destroy')}
              </Button>
            </>
          )}
        </div>
      </div>
      {(isRunning || isWaiting || isPending || isTerminating) && (
        <div className="h-1.5 bg-neutral-700 rounded-full overflow-hidden">
          {isRunning ? (
            <div className="h-full bg-yellow-400" style={{ width: `${progressWidth}%` }} />
          ) : (
            <div className={`h-full w-full ${isTerminating ? 'bg-orange-400/35' : 'bg-yellow-400/35'}`} />
          )}
        </div>
      )}
      {isRunning && challenge.instanceIP && (
        <div>
          <div className="flex items-center justify-between p-1 rounded-md">
            <span className="text-neutral-400 text-xs">{t('game.challengeModal.instance.address')}</span>
          </div>
          {challenge.instanceIP.map((ip, index) => (
            <div key={index} className="flex min-w-0 items-center justify-between gap-2 bg-neutral-900">
              <span className="min-w-0 break-all font-mono text-neutral-50 text-sm">{ip}</span>
              <Button
                variant="ghost"
                size="icon"
                className="shrink-0 !text-neutral-400 hover:!text-geek-400"
                aria-label={t('common.copyToClipboard', { defaultValue: 'Copy to clipboard' })}
                onClick={() => onCopy(ip)}
              >
                {isCopied[ip] ? <IconCheck size={16} /> : <IconCopy size={16} />}
              </Button>
            </div>
          ))}
        </div>
      )}
      {isRunning && accessKey && (
        <div className="space-y-2">
          <div className="flex items-center justify-between p-1 rounded-md">
            <span className="text-neutral-400 text-xs">{t('game.challengeModal.instance.sshAccess')}</span>
          </div>
          <div className="bg-neutral-900 px-2 py-1.5 font-mono text-neutral-50 text-sm break-all">
            {t('game.challengeModal.instance.sshUser')}: {accessUser}
          </div>
          {sshCommand && (
            <div className="flex min-w-0 items-center justify-between gap-2 bg-neutral-900">
              <span className="min-w-0 break-all font-mono text-neutral-50 text-sm">{sshCommand}</span>
              <Button
                variant="ghost"
                size="icon"
                className="shrink-0 !text-neutral-400 hover:!text-geek-400"
                aria-label={t('common.copyToClipboard', { defaultValue: 'Copy to clipboard' })}
                onClick={() => onCopy(sshCommand)}
              >
                {isCopied[sshCommand] ? <IconCheck size={16} /> : <IconCopy size={16} />}
              </Button>
            </div>
          )}
          <pre className="max-h-32 overflow-auto bg-neutral-900 p-2 font-mono text-xs text-neutral-300 break-all">
            {accessKey}
          </pre>
          <Button variant="secondary" size="sm" onClick={downloadKey}>
            <span className="flex items-center gap-2">
              <IconDownload size={16} />
              {t('game.challengeModal.instance.downloadKey')}
            </span>
          </Button>
        </div>
      )}
    </div>
  );
}
