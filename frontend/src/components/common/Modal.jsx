import { createPortal } from 'react-dom';
import { useEffect, useRef, useId } from 'react';
import { motion, AnimatePresence } from 'motion/react';
import { IconX } from '@tabler/icons-react';
import Button from './Button';
import { useModalPortal } from './ModalProvider';

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
  confirmText = 'CONFIRM',
  cancelText = 'CANCEL',
  onConfirm,
  confirmType = 'primary',
  className = '',
  showHeader = true,
  showCloseButton = true,
}) {
  const portalContainer = useModalPortal();
  const titleId = useId();
  const dialogRef = useRef(null);
  const triggerRef = useRef(null);
  const onCloseRef = useRef(onClose);
  useEffect(() => {
    onCloseRef.current = onClose;
  });

  // 保存触发元素，关闭时恢复焦点
  useEffect(() => {
    if (isOpen) {
      triggerRef.current = document.activeElement;
    }
  }, [isOpen]);

  // 焦点陷阱：将焦点移入模态框，并限制 Tab 键在内部循环
  useEffect(() => {
    if (!isOpen || !dialogRef.current) return;

    // 将焦点移至模态框内第一个可聚焦元素
    const focusableElements = dialogRef.current.querySelectorAll(FOCUSABLE_SELECTORS);
    if (focusableElements.length > 0) {
      focusableElements[0].focus();
    } else {
      dialogRef.current.focus();
    }

    const handleKeyDown = (e) => {
      if (e.key === 'Escape') {
        onCloseRef.current?.();
        return;
      }
      if (e.key !== 'Tab') return;

      const elements = dialogRef.current.querySelectorAll(FOCUSABLE_SELECTORS);
      if (elements.length === 0) return;

      const first = elements[0];
      const last = elements[elements.length - 1];

      if (e.shiftKey) {
        if (document.activeElement === first) {
          e.preventDefault();
          last.focus();
        }
      } else {
        if (document.activeElement === last) {
          e.preventDefault();
          first.focus();
        }
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      // 恢复焦点到触发元素
      triggerRef.current?.focus();
    };
  }, [isOpen]);

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
    : `relative w-full ${sizeClasses[size]} mx-4`;
  const dialogContainerClassName = isFullScreen
    ? `flex h-full flex-col overflow-hidden border border-neutral-300/30 bg-black/80 backdrop-blur-[8px] rounded-none ${className}`
    : `border border-neutral-300/30 rounded-lg bg-black/80 backdrop-blur-[8px] overflow-hidden ${className}`;
  const dialogBodyClassName = bodyClassName
    ? `${isFullScreen ? 'min-h-0 flex-1 ' : ''}${bodyClassName}`
    : isFullScreen
      ? 'p-6 min-h-0 flex-1 overflow-y-auto'
      : 'p-6 max-h-[70vh] overflow-y-auto';

  // Portal容器未就绪时不渲染
  if (!portalContainer) return null;

  // Confirm模式：简化的确认对话框
  if (variant === 'confirm') {
    return createPortal(
      <AnimatePresence>
        {isOpen && (
          <div className="fixed inset-0 z-[100] flex items-center justify-center">
            <motion.div
              className="fixed inset-0 bg-black/60 backdrop-blur-sm"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={onClose}
            />
            <motion.div
              ref={dialogRef}
              role="dialog"
              aria-modal="true"
              aria-labelledby={titleId}
              tabIndex={-1}
              className={`relative w-full ${sizeClasses.sm} m-4 p-6 border border-neutral-300 rounded-md bg-black/80 ${className}`}
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
            >
              <h3 id={titleId} className="text-lg font-mono text-neutral-50 mb-4">
                {title}
              </h3>
              <div className="text-neutral-300 mb-6">{children}</div>
              <div className="flex justify-end gap-4">
                <Button size="sm" variant="ghost" onClick={onClose}>
                  {cancelText}
                </Button>
                <Button size="sm" variant={confirmType === 'danger' ? 'danger' : 'primary'} onClick={onConfirm}>
                  {confirmText}
                </Button>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>,
      portalContainer
    );
  }

  // Default模式：完整的模态框
  return createPortal(
    <AnimatePresence>
      {isOpen && (
        <div className="fixed inset-0 z-[999] flex items-center justify-center">
          {/* 背景遮罩 */}
          <motion.div
            className="absolute inset-0 bg-black/60 backdrop-blur-sm"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
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
            initial={{ opacity: 0, y: 20, scale: 0.95 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 20, scale: 0.95 }}
            transition={{ duration: 0.2 }}
          >
            <div className={dialogContainerClassName}>
              {!showHeader && showCloseButton ? (
                <div className="pointer-events-none absolute right-4 top-4 z-10">
                  <Button
                    variant="ghost"
                    size="icon"
                    aria-label="Close dialog"
                    className="pointer-events-auto !bg-black/30 !text-neutral-300 backdrop-blur-sm hover:!text-neutral-100"
                    onClick={onClose}
                  >
                    <IconX size={18} />
                  </Button>
                </div>
              ) : null}
              {/* 头部 */}
              {showHeader ? <div className="p-6 border-b border-neutral-300/30">
                <div className="flex items-center justify-between">
                  <h2 id={titleId} className="text-xl font-mono text-neutral-50">
                    {title}
                  </h2>
                  {showCloseButton ? (
                    <Button
                      variant="ghost"
                      size="icon"
                      aria-label="Close dialog"
                      className="!bg-transparent !text-neutral-400 hover:!text-neutral-200"
                      onClick={onClose}
                    >
                      <IconX size={18} />
                    </Button>
                  ) : null}
                </div>
              </div> : null}

              {/* 内容区域 */}
              <div className={dialogBodyClassName}>{children}</div>

              {/* 底部 */}
              {footer && (
                <div className="p-6 border-t border-neutral-300/30 bg-black/40">
                  <div className="flex justify-end gap-3">{footer}</div>
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

// 子组件：模态框头部
Modal.Header = function ModalHeader({ children, className = '' }) {
  return <div className={`mb-4 ${className}`}>{children}</div>;
};

// 子组件：模态框主体
Modal.Body = function ModalBody({ children, className = '' }) {
  return <div className={className}>{children}</div>;
};

// 子组件：模态框底部
Modal.Footer = function ModalFooter({ children, className = '' }) {
  return <div className={`mt-4 ${className}`}>{children}</div>;
};

export default Modal;
