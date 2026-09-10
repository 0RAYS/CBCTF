import { IconBan } from '@tabler/icons-react';
import { Button, Modal } from '../../../common';

export default function VictimStopDialog({ t, isOpen, onClose, onConfirm, selectedCount, translationKey, contestId }) {
  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={t(`${translationKey}.modals.stopTitle`)}
      footer={
        <>
          <Button size="sm" variant="ghost" onClick={onClose}>
            {t('common.cancel')}
          </Button>
          <Button size="sm" variant="danger" onClick={onConfirm}>
            {t(`${translationKey}.modals.stopConfirm`)}
          </Button>
        </>
      }
    >
      <div className="flex items-center gap-3">
        {contestId && <IconBan size={20} className="text-red-400" />}
        <p className="text-neutral-300 font-mono">
          {t(`${translationKey}.modals.stopPrompt`, { count: selectedCount })}
        </p>
      </div>
    </Modal>
  );
}
