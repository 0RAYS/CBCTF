import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';

/**
 * 通用按钮组件
 * @param {Object} props
 * @param {React.ReactNode} props.children - 按钮内容
 * @param {React.ReactNode} props.icon - 按钮图标，显示在文字左侧
 * @param {function} props.onClick - 点击事件回调
 * @param {string} props.variant - 按钮风格: 'default', 'primary', 'danger', 'outline', 'ghost'
 * @param {string} props.size - 按钮大小: 'sm', 'md', 'lg', 'action', 'icon'
 * @param {string} props.align - 内容对齐方式: 'center', 'left', 'right', 'icon-left'
 * @param {string} props.textColor - 文字颜色，覆盖variant默认文字颜色
 * @param {boolean} props.fullWidth - 是否占满宽度
 * @param {boolean} props.disabled - 是否禁用
 * @param {boolean} props.loading - 是否加载中
 * @param {string} props.type - 按钮类型: 'button', 'submit', 'reset'
 * @param {string} props.className - 额外的自定义类名
 * @param {boolean} props.animate - 是否启用动画效果
 */
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

  // 按钮风格变体（已移除 hover glow，改用 border/opacity 过渡）
  const variants = {
    default: 'bg-black/30 border-neutral-300/30 text-neutral-400 hover:border-neutral-300 hover:text-neutral-300',
    primary: 'border-geek-400 text-geek-400 hover:bg-geek-400/10 hover:border-geek-400/80',
    danger: 'border-red-400 text-red-400 hover:bg-red-400/10 hover:border-red-400/80',
    outline:
      'border-neutral-300 text-neutral-300 hover:bg-neutral-300/10 hover:text-neutral-50 hover:border-neutral-50',
    ghost: 'border-transparent text-neutral-400 hover:text-neutral-200 hover:bg-white/5',
  };

  // 按钮尺寸
  const sizes = {
    sm: 'h-[36px] text-sm',
    md: 'h-[42px] text-sm',
    lg: 'h-[50px] text-lg tracking-wider',
    action: 'h-[40px] font-mono tracking-wider',
    icon: 'w-[36px] h-[36px] p-0', // 正方形图标按钮
  };

  // 内容对齐方式
  const alignments = {
    center: 'justify-center',
    left: 'justify-start',
    right: 'justify-end',
    'icon-left': 'justify-start',
  };

  // 内边距设置
  const getPadding = () => {
    if (size === 'icon') return '';
    if (align === 'icon-left') return 'px-4';
    return 'px-6';
  };

  // 基础样式（含 focus-visible ring，替代原先的 outline-none）
  const buttonClasses = `
    relative border rounded-md font-mono transition-colors
    focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-geek-400/70 focus-visible:ring-offset-1 focus-visible:ring-offset-black
    ${variants[variant] || variants.default}
    ${sizes[size] || sizes.md}
    ${getPadding()}
    ${fullWidth ? 'w-full' : 'inline-flex'}
    items-center ${alignments[align]}
    ${disabled || loading ? 'opacity-50 cursor-not-allowed' : ''}
    ${className}
  `;

  // 动画属性：使用 opacity 替代 scale，更克制
  const motionProps =
    animate && !disabled && !loading
      ? {
          whileHover: { opacity: 0.85 },
          whileTap: { opacity: 0.7 },
          transition: { duration: 0.15 },
        }
      : {};

  // 渲染按钮
  const ButtonComponent = animate ? motion.button : 'button';

  // 渲染内容
  const renderContent = () => {
    if (loading) {
      return (
        <span className="flex items-center justify-center">
          <svg
            className="animate-spin -ml-1 mr-2 h-4 w-4"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
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
      {...motionProps}
      {...rest}
    >
      {/* 主按钮 hover 填充动画（仅 primary 变体） */}
      {!disabled && !loading && variant === 'primary' && animate && size !== 'icon' && (
        <motion.div
          className="absolute inset-0 bg-geek-400/20"
          initial={{ scale: 0, opacity: 0 }}
          whileHover={{ scale: 1, opacity: 1 }}
          transition={{ duration: 0.2 }}
        />
      )}

      {/* 按钮内容 */}
      <span className="relative z-10 w-full">{renderContent()}</span>
    </ButtonComponent>
  );
}

export default Button;
