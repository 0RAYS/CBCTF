/**
 * 通用确认模态框（基于通用Modal）
 * @param {Object} props
 * @param {boolean} props.isOpen - 是否显示
 * @param {Function} props.onClose - 关闭回调
 * @param {Function} props.onConfirm - 确认回调
 * @param {string} props.title - 标题
 * @param {string} props.message - 提示信息
 * @param {string} [props.confirmText='CONFIRM'] - 确认按钮文本
 * @param {string} [props.type='default'] - 类型 (default/danger)
 */

import Modal from './Modal';
import ModalFooter from './ModalFooter';
import { useTranslation } from 'react-i18next';

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
