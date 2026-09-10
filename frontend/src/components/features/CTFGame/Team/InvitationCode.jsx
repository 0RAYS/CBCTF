import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Button } from '../../../common';

export default function InvitationCode({ code, isLeader, onCopyCode, onRefreshCode }) {
  const { t } = useTranslation();
  const [isCodeCopied, setIsCodeCopied] = useState(false);
  const timer = useRef(null);
  useEffect(() => () => clearTimeout(timer.current), []);

  return (
    <div className="space-y-2">
      <div className="text-neutral-400 text-sm">{t('game.team.settings.invitationCode')}</div>
      <div className="flex items-center gap-2">
        <div className="min-w-0 flex-1 flex items-center justify-between gap-2 p-3 bg-neutral-900 rounded-md">
          <span className="min-w-0 font-mono text-neutral-50 [overflow-wrap:anywhere]">{code}</span>
          <Button
            variant="ghost"
            size="icon"
            className="shrink-0"
            aria-label={isCodeCopied ? t('game.team.toast.inviteCopied') : t('common.copy', { defaultValue: 'Copy' })}
            onClick={() => {
              onCopyCode();
              clearTimeout(timer.current);
              setIsCodeCopied(true);
              timer.current = setTimeout(() => setIsCodeCopied(false), 2000);
            }}
          >
            {isCodeCopied ? '✓' : '📋'}
          </Button>
        </div>
        {isLeader && (
          <Button
            variant="outline"
            size="icon"
            className="shrink-0"
            aria-label={t('common.refresh')}
            onClick={onRefreshCode}
          >
            ↻
          </Button>
        )}
      </div>
    </div>
  );
}
