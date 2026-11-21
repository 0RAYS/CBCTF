import { useTranslation } from 'react-i18next';
import ModalFooter from '@/components/common/ModalFooter';

/** Standard footer for CRUD modals with cancel and action buttons. */
function CRUDModalFooter({ mode, onCancel, onSubmit, labels = {}, submitDisabled = false, submitLoading = false }) {
  const { t } = useTranslation();

  const defaultLabels = {
    cancel: t('common.cancel'),
    create: t('common.create'),
    save: t('common.save'),
    delete: t('common.delete'),
  };

  const mergedLabels = { ...defaultLabels, ...labels };

  const actionLabel = { create: mergedLabels.create, edit: mergedLabels.save, delete: mergedLabels.delete }[mode];
  if (!actionLabel) return null;

  return (
    <ModalFooter
      onCancel={onCancel}
      onSubmit={onSubmit}
      cancelLabel={mergedLabels.cancel}
      submitLabel={actionLabel}
      submitVariant={mode === 'delete' ? 'danger' : 'primary'}
      submitDisabled={submitDisabled}
      submitLoading={submitLoading}
    />
  );
}

export default CRUDModalFooter;
