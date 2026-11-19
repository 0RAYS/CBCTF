import { motion } from 'motion/react';

/** 通用卡片容器组件 */
function Card({ children, variant = 'default', padding = 'md', className = '', animate = false, onClick, ...rest }) {
  const variants = {
    default: 'bg-neutral-800/80 border border-neutral-600/60 rounded-md',
    bordered: 'bg-neutral-800/60 border border-neutral-600 rounded-md',
    glass: 'bg-neutral-700/40 border border-neutral-600/80 rounded-lg backdrop-blur-sm',
    dark: 'bg-neutral-900 border border-neutral-300/20 rounded-md',
  };

  const paddings = {
    none: '',
    sm: 'p-2',
    md: 'p-4',
    lg: 'p-6',
  };

  const cardClasses = `
    ${variants[variant] || variants.default}
    ${paddings[padding] ?? paddings.md}
    ${onClick ? 'cursor-pointer transition-colors hover:border-geek-400/50' : ''}
    ${className}
  `
    .trim()
    .replace(/\s+/g, ' ');

  const motionProps =
    animate && !onClick
      ? {
          initial: { opacity: 0, y: 20 },
          animate: { opacity: 1, y: 0 },
          transition: { duration: 0.3, ease: [0.25, 1, 0.5, 1] },
        }
      : animate && onClick
        ? {
            whileHover: { opacity: 0.85 },
            whileTap: { opacity: 0.7 },
            transition: { duration: 0.15 },
          }
        : {};

  const CardComponent = animate ? motion.div : 'div';

  return (
    <CardComponent className={cardClasses} onClick={onClick} {...motionProps} {...rest}>
      {children}
    </CardComponent>
  );
}

export default Card;
