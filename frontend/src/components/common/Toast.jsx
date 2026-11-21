import { forwardRef } from 'react';
import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import {
  IconAlertTriangle,
  IconBolt,
  IconCheck,
  IconCircleX,
  IconDots,
  IconInfoCircle,
  IconX,
} from '@tabler/icons-react';
import { toastVariants } from '@/config/motion';

const TOAST_CONFIG = {
  primary: {
    container: 'border-geek-400/50 bg-neutral-800/95 text-neutral-200',
    title: 'text-geek-300 font-mono',
    iconColor: 'text-geek-400',
    icon: IconBolt,
  },
  secondary: {
    container: 'border-geek-400/50 bg-neutral-800/95 text-neutral-200',
    title: 'text-geek-300 font-mono',
    iconColor: 'text-geek-400',
    icon: IconBolt,
  },
  success: {
    container: 'border-green-400/50 bg-neutral-800/95 text-neutral-200',
    title: 'text-green-300 font-mono',
    iconColor: 'text-green-400',
    icon: IconCheck,
  },
  warning: {
    container: 'border-amber-400/50 bg-neutral-800/95 text-neutral-200',
    title: 'text-amber-300 font-mono',
    iconColor: 'text-amber-400',
    icon: IconAlertTriangle,
  },
  danger: {
    container: 'border-red-400/50 bg-neutral-800/95 text-neutral-200',
    title: 'text-red-300 font-mono',
    iconColor: 'text-red-400',
    icon: IconX,
  },
  info: {
    container: 'border-geek-400/50 bg-neutral-800/95 text-neutral-200',
    title: 'text-geek-300 font-mono',
    iconColor: 'text-geek-400',
    icon: IconInfoCircle,
  },
  default: {
    container: 'border-neutral-500/50 bg-neutral-800/95 text-neutral-300',
    title: 'text-neutral-200 font-mono',
    iconColor: 'text-neutral-400',
    icon: IconDots,
  },
};

const ToastIcon = ({ color }) => {
  const config = TOAST_CONFIG[color] || TOAST_CONFIG.default;
  const Icon = config.icon;
  return <Icon size={20} stroke={2} className={config.iconColor} />;
};

const Toast = forwardRef(({ id, title, description, color = 'default', onClose, hasCloseButton = true }, ref) => {
  const { t } = useTranslation();
  const config = TOAST_CONFIG[color] || TOAST_CONFIG.default;

  return (
    <motion.div
      ref={ref}
      className={`border rounded-md ${config.container} p-3 shadow-lg max-w-md z-[10000] overflow-hidden`}
      variants={toastVariants}
      initial="hidden"
      animate="visible"
      exit="exit"
      layout
    >
      <div className="relative z-10 flex items-start space-x-3">
        <div className="flex-shrink-0 mt-0.5">
          <ToastIcon color={color} />
        </div>

        <div className="flex-1 min-w-0">
          {title && <h3 className={`font-medium text-sm [overflow-wrap:anywhere] ${config.title}`}>{title}</h3>}
          {description && <div className="mt-1 text-sm leading-relaxed [overflow-wrap:anywhere]">{description}</div>}
        </div>

        {hasCloseButton && (
          <button
            type="button"
            onClick={() => onClose && onClose(id)}
            className="pointer-events-auto flex-shrink-0 ml-1 text-current opacity-70 hover:opacity-100 transition-opacity focus-visible:ring-2 focus-visible:ring-geek-400 rounded-md p-2"
            aria-label={t('common.closeNotification')}
          >
            <IconCircleX size={16} stroke={2} />
          </button>
        )}
      </div>
    </motion.div>
  );
});

Toast.displayName = 'Toast';

export default Toast;
