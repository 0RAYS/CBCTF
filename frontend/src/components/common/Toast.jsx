import { forwardRef } from 'react';
import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import { toastVariants } from '../../config/motion';

const TOAST_CONFIG = {
  primary: {
    container: 'border-geek-400/70 bg-geek-900/70 text-geek-300 shadow-glow-primary',
    title: 'text-geek-300 font-mono',
    iconColor: 'text-geek-400',
    icon: <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />,
  },
  secondary: {
    container: 'border-geek-400/70 bg-geek-900/70 text-geek-300 shadow-glow-primary',
    title: 'text-geek-300 font-mono',
    iconColor: 'text-geek-400',
    icon: <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />,
  },
  success: {
    container: 'border-green-400/70 bg-green-900/70 text-green-300 shadow-glow-success',
    title: 'text-green-300 font-mono',
    iconColor: 'text-green-400',
    icon: <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />,
  },
  warning: {
    container: 'border-amber-400/70 bg-amber-900/70 text-amber-300 shadow-glow-warning',
    title: 'text-amber-300 font-mono',
    iconColor: 'text-amber-400',
    icon: (
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
      />
    ),
  },
  danger: {
    container: 'border-red-400/70 bg-red-900/70 text-red-300 shadow-glow-danger',
    title: 'text-red-300 font-mono',
    iconColor: 'text-red-400',
    icon: <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />,
  },
  info: {
    container: 'border-cyan-400/70 bg-cyan-900/70 text-cyan-300 shadow-glow-info',
    title: 'text-cyan-300 font-mono',
    iconColor: 'text-cyan-400',
    icon: (
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    ),
  },
  default: {
    container: 'border-neutral-500/50 bg-black/70 text-neutral-300 shadow-glow-muted',
    title: 'text-neutral-200 font-mono',
    iconColor: 'text-neutral-400',
    icon: (
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M8 12h.01M12 12h.01M16 12h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    ),
  },
};

const ToastIcon = ({ color }) => {
  const config = TOAST_CONFIG[color] || TOAST_CONFIG.default;
  return (
    <svg className={`w-5 h-5 ${config.iconColor}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      {config.icon}
    </svg>
  );
};

const Toast = forwardRef(({ id, title, description, color = 'default', onClose, hasCloseButton = true }, ref) => {
  const { t } = useTranslation();
  const config = TOAST_CONFIG[color] || TOAST_CONFIG.default;

  return (
    <motion.div
      ref={ref}
      className={`border rounded-md ${config.container} p-3.5 shadow-lg max-w-md z-[10000] overflow-hidden`}
      variants={toastVariants}
      initial="hidden"
      animate="visible"
      exit="exit"
      layout
    >
      {/* 添加一个微妙的动画背景 */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute inset-0 opacity-20 bg-gradient-to-r from-transparent via-white/5 to-transparent animate-pulse"></div>
      </div>

      <div className="relative z-10 flex items-start space-x-3">
        {/* 图标 */}
        <div className="flex-shrink-0 mt-0.5">
          <ToastIcon color={color} />
        </div>

        {/* 内容 */}
        <div className="flex-1 min-w-0">
          {title && <h3 className={`font-medium text-base ${config.title}`}>{title}</h3>}
          {description && (
            <div className="mt-1 text-sm opacity-90 font-light break-words overflow-wrap-anywhere">{description}</div>
          )}
        </div>

        {/* 关闭按钮 */}
        {hasCloseButton && (
          <button
            onClick={() => onClose && onClose(id)}
            className="flex-shrink-0 ml-1 text-current opacity-70 hover:opacity-100 transition-opacity focus:outline-none focus:ring-1 focus:ring-white/20 rounded-full p-1"
            aria-label={t('common.closeNotification')}
          >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="w-4 h-4">
              <path
                fillRule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z"
                clipRule="evenodd"
              />
            </svg>
          </button>
        )}
      </div>
    </motion.div>
  );
});

Toast.displayName = 'Toast';

export default Toast;
