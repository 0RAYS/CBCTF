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
import { useTranslation } from 'react-i18next';

function ConfirmModal({ isOpen, onClose, onConfirm, title, message, confirmText, type = 'default' }) {
  const { t } = useTranslation();

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={title}
      variant="confirm"
      confirmText={confirmText || t('common.confirm')}
      cancelText={t('common.cancel')}
      onConfirm={onConfirm}
      confirmType={type === 'danger' ? 'danger' : 'primary'}
    >
      {message}
    </Modal>
  );
}

export default ConfirmModal;
