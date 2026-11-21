import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getContestCheat } from '@/api/admin/contest';
import { toast } from '@/utils/toast';
import { Button, Modal, Spinner } from '@/components/common';
import {
  CheatEvidenceActors,
  CheatEvidenceSummary,
  CheatIp,
  CheatModels,
  CheatStatus,
  CheatVerdict,
} from '@/components/features/Admin/cheats/CheatEvidence';
import { formatDateTime } from '@/components/features/Admin/cheats/cheatFormat';
import CheatEvidencePanel from '@/components/features/Admin/cheats/CheatEvidencePanels';
import CheatEvidenceTimeline from '@/components/features/Admin/cheats/CheatEvidenceTimeline';

export default function CheatEvidenceDialog({ contestId, cheat, onClose, openUserDetail, openTeamDetail }) {
  const { t, i18n } = useTranslation();
  const [detail, setDetail] = useState(cheat);
  const [loading, setLoading] = useState(true);
  const language = i18n.language || 'en-US';

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    getContestCheat(contestId, cheat.id)
      .then((response) => {
        if (response.code !== 200) throw new Error(t('admin.contests.cheats.toast.fetchFailed'));
        if (!cancelled) setDetail(response.data);
      })
      .catch((error) => {
        if (!cancelled) toast.danger({ description: error.message || t('admin.contests.cheats.toast.fetchFailed') });
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [contestId, cheat.id, t]);

  const evidence = detail?.evidence;
  const fields = [
    ['eventId', detail.id],
    [
      'model',
      <CheatModels key="model" model={detail.model} openUserDetail={openUserDetail} openTeamDetail={openTeamDetail} />,
    ],
    ['ip', <CheatIp key="ip" ip={detail.ip} />],
    ['time', detail.time ? formatDateTime(detail.time, language) : '-'],
  ];
  return (
    <Modal
      isOpen
      onClose={onClose}
      title={t('admin.contests.cheats.modals.detailTitle')}
      size="lg"
      footer={
        <Button size="sm" variant="ghost" onClick={onClose}>
          {t('common.close')}
        </Button>
      }
    >
      <dl className="space-y-4">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {fields.map(([key, value]) => (
            <div key={key} className="min-w-0">
              <dt className="text-sm font-mono text-neutral-400 mb-1">{t(`admin.contests.cheats.detail.${key}`)}</dt>
              <dd className="text-neutral-300 break-words">{value}</dd>
            </div>
          ))}
        </div>
        <div>
          <dt className="text-sm font-mono text-neutral-400 mb-1">{t('admin.contests.cheats.detail.type')}</dt>
          <dd>
            <CheatVerdict type={detail.type} />
          </dd>
        </div>
        <div>
          <dt className="text-sm font-mono text-neutral-400 mb-1">{t('admin.contests.cheats.detail.reasonType')}</dt>
          <dd className="text-neutral-300">
            {detail.reason_type
              ? t(`admin.contests.cheats.reasonTypes.${detail.reason_type}`, detail.reason_type)
              : '-'}
          </dd>
        </div>
        <div>
          <dt className="text-sm font-mono text-neutral-400 mb-1">{t('admin.contests.cheats.detail.reason')}</dt>
          <dd className="text-neutral-300 bg-neutral-800 p-3 rounded whitespace-pre-wrap break-words">
            {detail.reason || '-'}
          </dd>
        </div>
        <div>
          <dt className="text-sm font-mono text-neutral-400 mb-1">{t('admin.contests.cheats.detail.status')}</dt>
          <dd className="flex flex-wrap items-center gap-3">
            <CheatStatus checked={detail.checked} />
            {detail.reviewed_by > 0 && (
              <span className="text-xs text-neutral-400 font-mono">
                {t('admin.contests.cheats.detail.reviewedBy')}: User-{detail.reviewed_by}
              </span>
            )}
            {detail.reviewed_at && (
              <span className="text-xs text-neutral-400 font-mono">
                {t('admin.contests.cheats.detail.reviewedAt')}: {formatDateTime(detail.reviewed_at, language)}
              </span>
            )}
          </dd>
        </div>
        {['comment', 'hash'].map(
          (key) =>
            detail[key] && (
              <div key={key}>
                <dt className="text-sm font-mono text-neutral-400 mb-1">{t(`admin.contests.cheats.detail.${key}`)}</dt>
                <dd className="text-neutral-300 font-mono text-sm bg-neutral-800 p-3 rounded whitespace-pre-wrap break-all">
                  {detail[key]}
                </dd>
              </div>
            )
        )}
        <div className="border-t border-neutral-300/10 pt-4 space-y-4">
          <dt className="text-sm font-mono text-neutral-400">{t('admin.contests.cheats.detail.evidence')}</dt>
          {loading && (
            <div className="flex items-center gap-2 text-neutral-400 text-sm">
              <Spinner size="sm" /> {t('admin.contests.cheats.evidence.loading')}
            </div>
          )}
          {!loading && evidence?.subject && (
            <div className="space-y-4">
              <CheatEvidenceSummary summary={detail.evidence_summary} language={language} />
              <CheatEvidenceActors
                actors={evidence.actors}
                openUserDetail={openUserDetail}
                openTeamDetail={openTeamDetail}
              />
              <CheatEvidencePanel evidence={evidence} language={language} />
              <div className="space-y-2">
                <div className="text-sm font-mono text-neutral-400">{t('admin.contests.cheats.evidence.timeline')}</div>
                <CheatEvidenceTimeline timeline={evidence.timeline} language={language} />
              </div>
            </div>
          )}
          {!loading && !evidence?.subject && (
            <div className="text-neutral-400 text-sm">{t('admin.contests.cheats.evidence.empty')}</div>
          )}
        </div>
      </dl>
    </Modal>
  );
}
