import { useTranslation } from 'react-i18next';
import { Button, Modal } from '../../../common';
import { CheatIp, CheatModels, CheatStatus, CheatVerdict } from './CheatEvidence';

export default function CheatEvidenceDialog({ cheat, onClose, openUserDetail, openTeamDetail }) {
  const { t, i18n } = useTranslation();
  const fields = [
    ['eventId', cheat.id],
    [
      'model',
      <CheatModels key="model" model={cheat.model} openUserDetail={openUserDetail} openTeamDetail={openTeamDetail} />,
    ],
    ['ip', <CheatIp key="ip" ip={cheat.ip} />],
    ['time', cheat.time ? new Date(cheat.time).toLocaleString(i18n.language || 'en-US') : '-'],
  ];
  return (
    <Modal
      isOpen
      onClose={onClose}
      title={t('admin.contests.cheats.modals.detailTitle')}
      size="lg"
      footer={
        <Button size="sm" variant="ghost" onClick={onClose}>
          {t('admin.contests.cheats.actions.close')}
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
            <CheatVerdict type={cheat.type} />
          </dd>
        </div>
        <div>
          <dt className="text-sm font-mono text-neutral-400 mb-1">{t('admin.contests.cheats.detail.reasonType')}</dt>
          <dd className="text-neutral-300">
            {cheat.reason_type ? t(`admin.contests.cheats.reasonTypes.${cheat.reason_type}`, cheat.reason_type) : '-'}
          </dd>
        </div>
        <div>
          <dt className="text-sm font-mono text-neutral-400 mb-1">{t('admin.contests.cheats.detail.reason')}</dt>
          <dd className="text-neutral-300 bg-neutral-800 p-3 rounded whitespace-pre-wrap break-words">
            {cheat.reason || '-'}
          </dd>
        </div>
        <div>
          <dt className="text-sm font-mono text-neutral-400 mb-1">{t('admin.contests.cheats.detail.status')}</dt>
          <dd>
            <CheatStatus checked={cheat.checked} />
          </dd>
        </div>
        {['comment', 'hash'].map(
          (key) =>
            cheat[key] && (
              <div key={key}>
                <dt className="text-sm font-mono text-neutral-400 mb-1">{t(`admin.contests.cheats.detail.${key}`)}</dt>
                <dd className="text-neutral-300 font-mono text-sm bg-neutral-800 p-3 rounded whitespace-pre-wrap break-all">
                  {cheat[key]}
                </dd>
              </div>
            )
        )}
      </dl>
    </Modal>
  );
}
