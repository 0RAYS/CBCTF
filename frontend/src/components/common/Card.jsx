import { motion } from 'motion/react';
import 'react';

/**
 * 通用卡片容器组件
 * @param {Object} props
 * @param {React.ReactNode} props.children - 卡片内容
 * @param {'default'|'bordered'|'glass'|'dark'} props.variant - 卡片风格
 * @param {'none'|'sm'|'md'|'lg'} props.padding - 内边距大小
 * @param {string} props.className - 额外的自定义类名
 * @param {boolean} props.animate - 是否启用动画效果
 * @param {function} props.onClick - 点击事件回调
 */
function Card({ children, variant = 'default', padding = 'md', className = '', animate = false, onClick, ...rest }) {
  // 卡片风格变体
  const variants = {
    default: 'bg-black/30 backdrop-blur-[2px] border border-neutral-300 rounded-md',
    bordered: 'bg-black/40 border border-neutral-300/30 rounded-md',
    glass: 'bg-black/20 backdrop-blur-[8px] border border-neutral-300/30 rounded-lg',
    dark: 'bg-neutral-900 border border-neutral-300 rounded-md',
  };

  // 内边距
  const paddings = {
    none: '',
    sm: 'p-2',
    md: 'p-4',
    lg: 'p-6',
  };

  // 基础样式
  const cardClasses = `
    ${variants[variant] || variants.default}
    ${paddings[padding] || paddings.md}
    ${onClick ? 'cursor-pointer transition-colors hover:border-geek-400/50' : ''}
    ${className}
  `
    .trim()
    .replace(/\s+/g, ' ');

  // 动画属性
  const motionProps =
    animate && !onClick
      ? {
          initial: { opacity: 0, y: 20 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.3 },
        }
      : animate && onClick
        ? {
            whileHover: { scale: 1.02, y: -2 },
            whileTap: { scale: 0.98 },
            transition: { duration: 0.2 },
          }
        : {};

  // 渲染卡片
  const CardComponent = animate ? motion.div : 'div';

  return (
    <CardComponent className={cardClasses} onClick={onClick} {...motionProps} {...rest}>
      {children}
    </CardComponent>
  );
}

// 子组件：卡片头部
Card.Header = function CardHeader({ children, className = '' }) {
  return <div className={`mb-4 ${className}`}>{children}</div>;
};

// 子组件：卡片主体
Card.Body = function CardBody({ children, className = '' }) {
  return <div className={className}>{children}</div>;
};

// 子组件：卡片底部
Card.Footer = function CardFooter({ children, className = '' }) {
  return <div className={`mt-4 ${className}`}>{children}</div>;
};

export default Card;
