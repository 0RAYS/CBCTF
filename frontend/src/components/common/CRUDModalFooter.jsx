import ModalButton from './ModalButton';
import { useTranslation } from 'react-i18next';

/**
 * Standard footer for CRUD modals with cancel and action buttons.
 * @param {Object} props
 * @param {string} props.mode - Current mode ('create' | 'edit' | 'delete')
 * @param {Function} props.onCancel - Cancel handler
 * @param {Function} props.onSubmit - Submit handler
 * @param {Object} [props.labels] - Override button labels { create, save, delete, cancel }
 */
function CRUDModalFooter({ mode, onCancel, onSubmit, labels = {} }) {
  const { t } = useTranslation();

  const defaultLabels = {
    cancel: t('common.cancel'),
    create: t('common.create'),
    save: t('common.save'),
    delete: t('common.delete'),
  };

  const mergedLabels = { ...defaultLabels, ...labels };

  const actionLabel =
    mode === 'create' ? mergedLabels.create : mode === 'edit' ? mergedLabels.save : mergedLabels.delete;

  return (
    <>
      <ModalButton variant="default" onClick={onCancel}>
        {mergedLabels.cancel}
      </ModalButton>
      <ModalButton variant={mode === 'delete' ? 'danger' : 'primary'} onClick={onSubmit}>
        {actionLabel}
      </ModalButton>
    </>
  );
}

export default CRUDModalFooter;
