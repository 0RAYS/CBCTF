import { useTranslation } from 'react-i18next';
import { Button, Modal } from '@/components/common';

export default function ChallengeCopyDialog({
  isOpen = false,
  challenge,
  contests = [],
  contestId = 0,
  onContestChange,
  onClose,
  onConfirm,
}) {
  const { t } = useTranslation();
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      size="md"
      title={t('admin.challenge.copy.title')}
      footer={
        <>
          <Button variant="ghost" size="sm" onClick={onClose}>
            {t('common.cancel')}
          </Button>
          <Button variant="primary" size="sm" onClick={() => onConfirm(contestId)}>
            {t('admin.challenge.copy.confirm')}
          </Button>
        </>
      }
    >
      <p className="text-sm text-neutral-300 mb-4">
        <span className="font-mono text-neutral-50">{challenge?.name}</span>
        {' — '}
        {t('admin.challenge.copy.description')}
      </p>
      <label className="block text-sm font-mono text-neutral-400 mb-1">
        {t('admin.challengeModal.labels.contest')}
      </label>
      <select
        value={String(contestId || 0)}
        onChange={(e) => onContestChange(Number(e.target.value))}
        className="select-custom select-custom-lg w-full"
      >
        <option value="0">{t('admin.challengeModal.contest.unassigned')}</option>
        {contests.map((contest) => (
          <option key={contest.id} value={contest.id}>
            {contest.name}
          </option>
        ))}
      </select>
      <p className="text-xs font-mono text-neutral-500 mt-3">{t('admin.challenge.copy.notice')}</p>
    </Modal>
  );
}
