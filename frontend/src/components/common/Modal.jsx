import { createPortal } from 'react-dom';
import { useEffect, useRef, useId } from 'react';
import { motion, AnimatePresence } from 'motion/react';
import { IconX } from '@tabler/icons-react';
import Button from './Button';
import { useModalPortal } from './ModalProvider';
import { backdropVariants, panelVariants } from '../../config/motion';
import { useTranslation } from 'react-i18next';

// Only the topmost dialog handles keyboard input; nested dialogs share the scroll lock.
const openDialogs = [];
let savedBodyOverflow;
let dialogOrder = 0;

// 可聚焦元素选择器
const FOCUSABLE_SELECTORS = [
  'a[href]',
  'button:not([disabled])',
  'textarea:not([disabled])',
  'input:not([disabled])',
  'select:not([disabled])',
  '[tabindex]:not([tabindex="-1"])',
].join(', ');

/**
 * 统一模态框组件（合并ConfirmModal和AdminModal）
 * @param {Object} props
 * @param {boolean} props.isOpen - 是否显示模态框
 * @param {function} props.onClose - 关闭模态框的回调函数
 * @param {string} props.title - 模态框标题
 * @param {React.ReactNode} props.children - 模态框内容
 * @param {React.ReactNode} props.footer - 模态框底部内容（default模式使用）
 * @param {'sm'|'md'|'lg'|'xl'|'2xl'} props.size - 模态框大小
 * @param {'default'|'confirm'} props.variant - 模态框变体
 * @param {string} props.confirmText - 确认按钮文本（confirm模式使用）
 * @param {string} props.cancelText - 取消按钮文本（confirm模式使用）
 * @param {function} props.onConfirm - 确认回调（confirm模式使用）
 * @param {'primary'|'danger'} props.confirmType - 确认按钮类型（confirm模式使用）
 * @param {string} props.className - 额外的自定义类名
 */
function Modal({
  isOpen,
  onClose,
  title,
  children,
  footer,
  bodyClassName = '',
  size = 'md',
  variant = 'default',
  confirmText,
  cancelText,
  onConfirm,
  confirmType = 'primary',
  className = '',
  showHeader = true,
  showCloseButton = true,
}) {
  const { t } = useTranslation();
  const portalContainer = useModalPortal();
  const titleId = useId();
  const dialogRef = useRef(null);
  const triggerRef = useRef(null);
  const onCloseRef = useRef(onClose);
  useEffect(() => {
    onCloseRef.current = onClose;
  });

  useEffect(() => {
    if (!isOpen || !portalContainer || !dialogRef.current) return;

    const dialog = dialogRef.current;
    triggerRef.current = document.activeElement;
    if (openDialogs.length === 0) {
      savedBodyOverflow = document.body.style.overflow;
      document.body.style.overflow = 'hidden';
    }
    openDialogs.push(dialog);
    dialog.parentElement.style.zIndex = String(1000 + dialogOrder++);

    const getFocusableElements = () =>
      [...dialog.querySelectorAll(FOCUSABLE_SELECTORS)].filter((element) => element.getClientRects().length > 0);
    (getFocusableElements()[0] || dialog).focus({ preventScroll: true });

    const handleKeyDown = (e) => {
      if (openDialogs.at(-1) !== dialog) return;
      if (e.key === 'Escape') {
        e.preventDefault();
        e.stopPropagation();
        onCloseRef.current?.();
        return;
      }
      if (e.key !== 'Tab') return;

      const elements = getFocusableElements();
      const first = elements[0] || dialog;
      const last = elements.at(-1) || dialog;
      if (
        !dialog.contains(document.activeElement) ||
        document.activeElement === dialog ||
        (e.shiftKey ? document.activeElement === first : document.activeElement === last)
      ) {
        e.preventDefault();
        (e.shiftKey ? last : first).focus();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      const wasTopmost = openDialogs.at(-1) === dialog;
      openDialogs.splice(openDialogs.indexOf(dialog), 1);
      if (openDialogs.length === 0) {
        document.body.style.overflow = savedBodyOverflow;
        dialogOrder = 0;
      }
      if (wasTopmost) {
        const target = triggerRef.current?.isConnected ? triggerRef.current : openDialogs.at(-1);
        target?.focus({ preventScroll: true });
      }
    };
  }, [isOpen, portalContainer]);

  // 根据size确定宽度
  const sizeClasses = {
    sm: 'max-w-[400px]',
    md: 'max-w-[600px]',
    lg: 'max-w-[800px]',
    xl: 'max-w-[1000px]',
    '2xl': 'max-w-[1200px]',
  };
  const isFullScreen = size === 'full';
  const dialogWrapperClassName = isFullScreen
    ? 'relative h-[100dvh] w-[100vw] max-w-none'
    : `relative min-w-0 w-full ${sizeClasses[size]} mx-4`;
  const dialogContainerClassName = isFullScreen
    ? `flex h-full flex-col overflow-hidden border border-neutral-600/50 bg-neutral-900/90 backdrop-blur-[8px] rounded-none ${className}`
    : `flex max-h-[calc(100dvh-2rem)] flex-col border border-neutral-600/70 rounded-lg bg-neutral-900/95 shadow-2xl backdrop-blur-[8px] overflow-hidden ${className}`;
  const dialogBodyClassName = `min-h-0 min-w-0 flex-1 overflow-y-auto overscroll-contain ${bodyClassName || 'p-4 sm:p-6'}`;

  // Portal容器未就绪时不渲染
  if (!portalContainer) return null;

  // Confirm模式: 简化的确认对话框
  if (variant === 'confirm') {
    return createPortal(
      <AnimatePresence>
        {isOpen && (
          <div className="fixed inset-0 z-[100] flex items-center justify-center">
            <motion.div
              className="fixed inset-0 bg-neutral-900/70 backdrop-blur-sm"
              variants={backdropVariants}
              initial="hidden"
              animate="visible"
              exit="exit"
              onClick={onClose}
            />
            <motion.div
              ref={dialogRef}
              role="dialog"
              aria-modal="true"
              aria-labelledby={titleId}
              tabIndex={-1}
              className={`relative w-full ${sizeClasses.sm} max-h-[calc(100dvh-2rem)] overflow-y-auto overscroll-contain m-4 p-6 border border-neutral-600/60 rounded-md bg-neutral-800/95 shadow-2xl ${className}`}
              variants={panelVariants}
              initial="hidden"
              animate="visible"
              exit="exit"
            >
              <h3 id={titleId} className="text-lg font-mono text-neutral-50 mb-4">
                {title}
              </h3>
              <div className="text-neutral-300 mb-6">{children}</div>
              <div className="flex flex-wrap justify-end gap-3">
                <Button size="sm" variant="ghost" onClick={onClose}>
                  {cancelText ?? t('common.cancel')}
                </Button>
                <Button size="sm" variant={confirmType === 'danger' ? 'danger' : 'primary'} onClick={onConfirm}>
                  {confirmText ?? t('common.confirm')}
                </Button>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>,
      portalContainer
    );
  }

  // Default模式: 完整的模态框
  return createPortal(
    <AnimatePresence>
      {isOpen && (
        <div className="fixed inset-0 z-[999] flex items-center justify-center">
          {/* 背景遮罩 */}
          <motion.div
            className="absolute inset-0 bg-neutral-900/70 backdrop-blur-sm"
            variants={backdropVariants}
            initial="hidden"
            animate="visible"
            exit="exit"
            onClick={onClose}
          />

          {/* 模态框主体 */}
          <motion.div
            ref={dialogRef}
            role="dialog"
            aria-modal="true"
            aria-labelledby={titleId}
            tabIndex={-1}
            className={dialogWrapperClassName}
            variants={panelVariants}
            initial="hidden"
            animate="visible"
            exit="exit"
          >
            <div className={dialogContainerClassName}>
              {!showHeader && (
                <h2 id={titleId} className="sr-only">
                  {title}
                </h2>
              )}
              {!showHeader && showCloseButton ? (
                <div className="pointer-events-none absolute right-4 top-4 z-10">
                  <Button
                    variant="ghost"
                    size="icon"
                    aria-label={t('common.close')}
                    className="pointer-events-auto !bg-black/30 !text-neutral-300 backdrop-blur-sm hover:!text-neutral-100"
                    onClick={onClose}
                  >
                    <IconX size={18} />
                  </Button>
                </div>
              ) : null}
              {/* 头部 */}
              {showHeader ? (
                <div className="shrink-0 p-4 sm:p-6 border-b border-neutral-600/50">
                  <div className="flex items-start justify-between gap-3">
                    <h2
                      id={titleId}
                      className="min-w-0 text-lg sm:text-xl font-mono text-neutral-50 [overflow-wrap:anywhere]"
                    >
                      {title}
                    </h2>
                    {showCloseButton ? (
                      <Button
                        variant="ghost"
                        size="icon"
                        aria-label={t('common.close')}
                        className="shrink-0 !bg-transparent !text-neutral-400 hover:!text-neutral-200"
                        onClick={onClose}
                      >
                        <IconX size={18} />
                      </Button>
                    ) : null}
                  </div>
                </div>
              ) : null}

              {/* 内容区域 */}
              <div className={dialogBodyClassName}>{children}</div>

              {/* 底部 */}
              {footer && (
                <div className="shrink-0 p-4 sm:p-6 border-t border-neutral-600/50 bg-neutral-800/40">
                  <div className="flex flex-wrap justify-end gap-3">{footer}</div>
                </div>
              )}
            </div>
          </motion.div>
        </div>
      )}
    </AnimatePresence>,
    portalContainer
  );
}

// 子组件: 模态框头部
Modal.Header = function ModalHeader({ children, className = '' }) {
  return <div className={`mb-4 ${className}`}>{children}</div>;
};

// 子组件: 模态框主体
Modal.Body = function ModalBody({ children, className = '' }) {
  return <div className={className}>{children}</div>;
};

// 子组件: 模态框底部
Modal.Footer = function ModalFooter({ children, className = '' }) {
  return <div className={`mt-4 ${className}`}>{children}</div>;
};

export default Modal;
