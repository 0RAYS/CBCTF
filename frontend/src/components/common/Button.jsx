import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import Spinner from './Spinner';

/** 通用按钮组件 */
function Button({
  children,
  icon,
  onClick,
  variant = 'default',
  size = 'md',
  align = 'center',
  textColor = '',
  fullWidth = false,
  disabled = false,
  loading = false,
  type = 'button',
  className = '',
  animate = true,
  ...rest
}) {
  const { t } = useTranslation();

  const variants = {
    default: 'bg-neutral-800/40 border-neutral-600/50 text-neutral-400 hover:border-neutral-400 hover:text-neutral-300',
    primary: 'border-geek-400 text-geek-400 hover:bg-geek-400/10 hover:border-geek-400/80',
    danger: 'border-red-400 text-red-400 hover:bg-red-400/10 hover:border-red-400/80',
    outline:
      'border-neutral-300 text-neutral-300 hover:bg-neutral-300/10 hover:text-neutral-50 hover:border-neutral-50',
    ghost: 'border-transparent text-neutral-400 hover:text-neutral-200 hover:bg-white/5',
  };

  const sizes = {
    sm: 'h-[40px] text-sm',
    md: 'h-[42px] text-sm',
    lg: 'h-[50px] text-lg tracking-wider',
    action: 'h-[40px] font-mono tracking-wider',
    icon: 'w-[40px] h-[40px] p-0', // min 40px — dense admin UI compromise
  };

  const alignments = {
    center: 'justify-center',
    left: 'justify-start',
    right: 'justify-end',
    'icon-left': 'justify-start',
  };

  const getPadding = () => {
    if (size === 'icon') return '';
    if (align === 'icon-left') return 'px-4';
    return 'px-6';
  };

  const buttonClasses = `
    relative border rounded-md font-mono transition-colors
    focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-geek-400/70 focus-visible:ring-offset-1 focus-visible:ring-offset-black
    ${variants[variant] || variants.default}
    ${sizes[size] || sizes.md}
    ${getPadding()}
    ${fullWidth ? 'w-full flex' : 'inline-flex'}
    items-center ${alignments[align]}
    ${disabled || loading ? 'opacity-50 cursor-not-allowed' : ''}
    ${className}
  `;

  const motionProps =
    animate && !disabled && !loading
      ? {
          whileHover: { opacity: 0.85 },
          whileTap: { opacity: 0.7 },
          transition: { duration: 0.15 },
        }
      : {};

  const ButtonComponent = animate ? motion.button : 'button';

  const renderContent = () => {
    if (loading) {
      return (
        <span className="flex items-center justify-center">
          <Spinner size="sm" border="md" colorClassName="border-current" className="mr-2" />
          {t('common.loading')}
        </span>
      );
    }

    if (size === 'icon') {
      return <div className="flex items-center justify-center w-full h-full">{children}</div>;
    }

    if (icon || align === 'icon-left') {
      return (
        <div className="flex items-center gap-2 whitespace-nowrap">
          {icon}
          <span className={textColor || ''}>{children}</span>
        </div>
      );
    }

    return <span className={textColor || ''}>{children}</span>;
  };

  return (
    <ButtonComponent
      type={type}
      className={buttonClasses}
      onClick={disabled || loading ? undefined : onClick}
      disabled={disabled || loading}
      aria-busy={loading || undefined}
      {...motionProps}
      {...rest}
    >
      <span className="relative z-10 w-full">{renderContent()}</span>
    </ButtonComponent>
  );
}

export default Button;
