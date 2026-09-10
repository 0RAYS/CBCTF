import Notices from '../../components/features/CTFGame/Notice/Notice';
import { useState, useEffect } from 'react';
import { getContestNotices } from '../../api/contest';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Button, Card } from '../../components/common';

function GameNoticePage() {
  const { contestId } = useParams();
  const { t, i18n } = useTranslation();
  const [request, setRequest] = useState({ contestId: null, status: 'loading', notices: [] });
  const [attempt, setAttempt] = useState(0);

  useEffect(() => {
    let active = true;
    setRequest({ contestId, status: 'loading', notices: [] });
    getContestNotices(contestId)
      .then((res) => {
        if (!active) return;
        if (res.code !== 200 || !res.data) throw new Error('Failed to load notices');
        const notices = res.data.notices ?? [];
        if (!Array.isArray(notices)) throw new Error('Invalid notices response');
        setRequest({ contestId, status: 'success', notices });
      })
      .catch(() => {
        if (active) setRequest({ contestId, status: 'error', notices: [] });
      });
    // Ignore responses from a previous contest, retry, or unmounted page.
    return () => {
      active = false;
    };
  }, [contestId, attempt]);

  const status = request.contestId === contestId ? request.status : 'loading';

  if (status !== 'success' || request.notices.length === 0) {
    return (
      <div className="contest-container mx-auto">
        <Card className="space-y-4" aria-busy={status === 'loading'}>
          <h1 className="font-mono text-lg text-neutral-50">{t('nav.notice')}</h1>
          {status === 'error' ? (
            <>
              <p role="alert" className="text-neutral-300">
                {t('admin.contests.notices.toast.fetchFailed')}
              </p>
              <Button
                variant="primary"
                size="sm"
                onClick={() => {
                  setRequest({ contestId, status: 'loading', notices: [] });
                  setAttempt((current) => current + 1);
                }}
              >
                {t('common.retry', { defaultValue: 'Retry' })}
              </Button>
            </>
          ) : (
            <p role="status" className="text-neutral-400">
              {status === 'loading' ? t('common.loading') : t('common.noData')}
            </p>
          )}
        </Card>
      </div>
    );
  }

  return (
    <Notices
      key={contestId}
      notices={request.notices.map((notice) => ({
        ...notice,
        type: notice.type || 'normal',
        timestamp: new Date(notice.created_at).toLocaleString(i18n.language),
      }))}
    />
  );
}

export default GameNoticePage;
