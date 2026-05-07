import { IconBan, IconPlayerPlay, IconSearch } from '@tabler/icons-react';
import { Chip, Modal, Pagination } from '../../../../components/common';
import ModalButton from '../../../../components/common/ModalButton';

export function StartContainerModal({
  t,
  isOpen,
  onClose,
  onConfirm,
  selectedChallenges,
  challenges,
  randomTeamPercentage,
  selectedTeamCount,
  totalTeamCount,
  victimDurationSeconds,
  formatVictimDuration,
}) {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={t('admin.contests.containers.modals.startTitle')}
      footer={
        <>
          <ModalButton onClick={onClose}>{t('common.cancel')}</ModalButton>
          <ModalButton variant="primary" onClick={onConfirm}>
            {t('admin.contests.containers.modals.startConfirm')}
          </ModalButton>
        </>
      }
    >
      <div className="space-y-4">
        <div className="flex items-center gap-3">
          <IconPlayerPlay size={20} className="text-geek-400" />
          <p className="text-neutral-300 font-mono">{t('admin.contests.containers.modals.startPrompt')}</p>
        </div>

        <div className="space-y-4">
          <div>
            <h4 className="text-sm font-mono text-neutral-400 mb-2">
              {t('admin.contests.containers.modals.selectedChallenges', { count: selectedChallenges.length })}
            </h4>
            <div className="max-h-32 overflow-y-auto border border-neutral-300/30 rounded-md bg-black/10 p-2">
              {selectedChallenges.map((challengeId) => {
                const challenge = challenges.find((item) => item.id === challengeId);
                return challenge ? (
                  <div key={challengeId} className="text-sm font-mono text-geek-400 py-1">
                    - {challenge.name}
                  </div>
                ) : null;
              })}
            </div>
          </div>

          <InfoBox title={t('admin.contests.containers.modals.teamRatioTitle')}>
            <p className="text-sm font-mono text-geek-400">
              {t('admin.contests.containers.modals.teamRatioValue', { ratio: randomTeamPercentage })}
            </p>
            <p className="text-xs font-mono text-neutral-400">
              {t('admin.contests.containers.modals.teamRatioHint', { count: selectedTeamCount, total: totalTeamCount })}
            </p>
          </InfoBox>

          <InfoBox title={t('admin.contests.containers.modals.durationTitle')}>
            <p className="text-sm font-mono text-geek-400">
              {t('admin.contests.containers.modals.durationValue', { seconds: victimDurationSeconds })}
            </p>
            <p className="text-xs font-mono text-neutral-400">
              {t('admin.contests.containers.modals.durationHint', {
                value: formatVictimDuration(victimDurationSeconds),
              })}
            </p>
          </InfoBox>
        </div>

        <div className="bg-neutral-800/50 border border-neutral-600/30 rounded-md p-3">
          <p className="text-xs font-mono text-neutral-400">
            {t('admin.contests.containers.modals.summaryPrefix')}
            <span className="text-geek-400">{selectedChallenges.length}</span>
            {t('admin.contests.containers.modals.summaryMiddle')}
            <span className="text-geek-400">{selectedTeamCount}</span>
            {t('admin.contests.containers.modals.summaryEquals')}
            <span className="text-green-400"> {selectedChallenges.length * selectedTeamCount}</span>
            {t('admin.contests.containers.modals.summarySuffix')}
          </p>
        </div>

        <div className="border border-amber-400/40 rounded-md bg-amber-400/10 p-3">
          <p className="text-xs font-mono text-amber-200">{t('admin.contests.containers.modals.startWarning')}</p>
        </div>
      </div>
    </Modal>
  );
}

export function ChallengeDetailsModal({
  t,
  isOpen,
  onClose,
  detailChallenges,
  detailChallengeTotal,
  detailChallengePage,
  challengePageSize,
  challengeSearch,
  selectedChallenges,
  typeLabels,
  onChallengeSearchChange,
  onChallengeSelectionChange,
  onPageChange,
  getChallengeCategoryChipClass,
  getChallengeTypeChipClass,
}) {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={t('admin.contests.containers.modals.challengeDetailsTitle')}
      size="xl"
      footer={<ModalButton onClick={onClose}>{t('common.cancel')}</ModalButton>}
    >
      <div className="space-y-4">
        <div className="relative">
          <IconSearch
            size={14}
            className="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-500 pointer-events-none"
          />
          <input
            type="text"
            value={challengeSearch}
            onChange={(e) => onChallengeSearchChange(e.target.value)}
            placeholder={t('admin.contests.containers.quickActions.searchPlaceholder')}
            className="w-full h-10 pl-10 pr-3 bg-black/20 border border-neutral-300/30 rounded-md text-sm text-neutral-50 placeholder-neutral-500 focus:outline-none focus:border-geek-400 transition-all duration-200"
          />
        </div>

        <p className="text-sm font-mono text-neutral-400">
          {t('admin.contests.containers.modals.challengeDetailsHint', { count: detailChallengeTotal })}
        </p>

        {detailChallenges.length === 0 ? (
          <div className="border border-neutral-300/20 rounded-md bg-black/10 p-4 text-sm font-mono text-neutral-500">
            {t('admin.contests.containers.modals.challengeDetailsEmpty')}
          </div>
        ) : (
          <div className="space-y-4">
            {detailChallenges.map((challenge) => (
              <ChallengeDetailRow
                key={challenge.id}
                t={t}
                challenge={challenge}
                selected={selectedChallenges.includes(challenge.id)}
                typeLabels={typeLabels}
                onChange={(checked) => onChallengeSelectionChange(challenge.id, checked)}
                getChallengeCategoryChipClass={getChallengeCategoryChipClass}
                getChallengeTypeChipClass={getChallengeTypeChipClass}
              />
            ))}
          </div>
        )}

        {Math.ceil(detailChallengeTotal / challengePageSize) > 1 ? (
          <div className="pt-2 border-t border-neutral-300/20">
            <Pagination
              total={Math.ceil(detailChallengeTotal / challengePageSize)}
              current={detailChallengePage}
              pageSize={challengePageSize}
              onChange={onPageChange}
              showTotal
              totalItems={detailChallengeTotal}
            />
          </div>
        ) : null}
      </div>
    </Modal>
  );
}

export function StopContainerModal({ t, isOpen, onClose, onConfirm, selectedCount }) {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={t('admin.contests.containers.modals.stopTitle')}
      footer={
        <>
          <ModalButton onClick={onClose}>{t('common.cancel')}</ModalButton>
          <ModalButton variant="danger" onClick={onConfirm}>
            {t('admin.contests.containers.modals.stopConfirm')}
          </ModalButton>
        </>
      }
    >
      <div className="flex items-center gap-3">
        <IconBan size={20} className="text-red-400" />
        <p className="text-neutral-300 font-mono">
          {t('admin.contests.containers.modals.stopPrompt', { count: selectedCount })}
        </p>
      </div>
    </Modal>
  );
}

function InfoBox({ title, children }) {
  return (
    <div>
      <h4 className="text-sm font-mono text-neutral-400 mb-2">{title}</h4>
      <div className="border border-neutral-300/30 rounded-md bg-black/10 p-3 space-y-2">{children}</div>
    </div>
  );
}

function ChallengeDetailRow({
  t,
  challenge,
  selected,
  typeLabels,
  onChange,
  getChallengeCategoryChipClass,
  getChallengeTypeChipClass,
}) {
  return (
    <div className="border border-neutral-300/20 rounded-md bg-black/10 p-4 space-y-4">
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0 space-y-2">
          <div className="flex items-center gap-2 flex-wrap">
            <h4 className="text-base font-mono text-neutral-50 break-all">{challenge.name}</h4>
            {challenge.category ? (
              <Chip
                size="sm"
                label={challenge.category}
                colorClass={getChallengeCategoryChipClass(challenge.category)}
              />
            ) : null}
            {challenge.type ? (
              <Chip
                size="sm"
                label={typeLabels[challenge.type] || challenge.type}
                colorClass={getChallengeTypeChipClass(challenge.type)}
              />
            ) : null}
            {challenge.hidden ? (
              <Chip size="sm" label={t('admin.contests.challenges.hidden')} colorClass="bg-red-400/20 text-red-400" />
            ) : null}
          </div>
          <div className="flex items-center gap-3 flex-wrap text-xs font-mono text-neutral-400">
            <span>{t('admin.contests.containers.modals.challengeScore', { score: challenge.score || 0 })}</span>
            <span>{t('admin.contests.containers.modals.challengeSolvers', { count: challenge.solvers || 0 })}</span>
            <span>{t('admin.contests.containers.modals.challengeAttempts', { count: challenge.attempt || 0 })}</span>
            <span>{t('admin.contests.containers.modals.challengeId', { id: challenge.id })}</span>
          </div>
        </div>

        <label className="inline-flex items-center gap-2 text-sm font-mono text-neutral-300 cursor-pointer shrink-0">
          <input
            type="checkbox"
            checked={selected}
            onChange={(e) => onChange(e.target.checked)}
            className="w-4 h-4 rounded border-neutral-300/30 text-geek-400 focus:ring-geek-400 focus:ring-offset-0 bg-black/20"
          />
          {t('admin.contests.containers.modals.selectChallenge')}
        </label>
      </div>
    </div>
  );
}
