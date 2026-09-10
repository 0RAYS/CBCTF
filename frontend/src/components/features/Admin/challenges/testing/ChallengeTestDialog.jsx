import { IconLoader2 } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import Modal from '../../../../common/Modal';
import MarkdownContent from '../../../../common/MarkdownContent';
import TestAttachments from './TestAttachments';
import TestInstancePanel from './TestInstancePanel';
import useTestSession from './useTestSession';

function OpenTestDialog({ challenge, onClose }) {
  const { t } = useTranslation();
  const session = useTestSession(challenge?.id);

  return (
    <Modal
      isOpen
      onClose={onClose}
      size="lg"
      className="!bg-black/80 !border-neutral-300 !rounded-md"
      bodyClassName="p-5 space-y-5"
      title={
        <span className="flex flex-wrap items-center gap-x-4 gap-y-2">
          <span className="text-geek-400 text-base">{t('admin.challenge.testModal.title')}</span>
          <span className="text-2xl text-neutral-50">{challenge?.name}</span>
          <span className="text-yellow-400 text-base">{challenge?.type}</span>
        </span>
      }
    >
      <section className="space-y-1.5 min-w-0">
        <h3 className="text-neutral-400 font-mono text-sm">{t('admin.challenge.testModal.sections.description')}</h3>
        <MarkdownContent className="text-neutral-50 prose-invert prose-sm max-w-none">
          {challenge?.description || ''}
        </MarkdownContent>
      </section>
      {session.loading.status && (
        <div role="status" className="flex items-center gap-2 text-neutral-400 text-sm">
          <IconLoader2 size={16} className="animate-spin" aria-hidden="true" />
          {t('common.loading')}
        </div>
      )}
      <TestAttachments challenge={challenge} file={session.testStatus?.file} />
      {challenge?.type === 'pods' && <TestInstancePanel {...session} />}
    </Modal>
  );
}

export default function ChallengeTestDialog({ challenge, isOpen, onClose }) {
  // Closing unmounts the session; changing the challenge creates an independent one.
  return isOpen ? <OpenTestDialog key={challenge?.id} challenge={challenge} onClose={onClose} /> : null;
}
