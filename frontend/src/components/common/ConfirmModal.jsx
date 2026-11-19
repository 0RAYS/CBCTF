import Modal from './Modal';
import ModalFooter from './ModalFooter';
import { useTranslation } from 'react-i18next';

/** 通用确认模态框（基于通用Modal） */
function ConfirmModal({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  children,
  confirmText,
  cancelText,
  type = 'default',
  loading = false,
  disabled = false,
}) {
  const { t } = useTranslation();
  const close = () => {
    if (!loading) onClose?.();
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={close}
      title={title}
      size="sm"
      showCloseButton={false}
      footer={
        <ModalFooter
          onCancel={close}
          onSubmit={onConfirm}
          cancelLabel={cancelText ?? t('common.cancel')}
          submitLabel={confirmText ?? t('common.confirm')}
          submitVariant={type === 'danger' ? 'danger' : 'primary'}
          submitLoading={loading}
          submitDisabled={disabled}
        />
      }
    >
      {children ?? message}
    </Modal>
  );
}

export default ConfirmModal;
