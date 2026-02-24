import 'react';

const sizeClasses = {
  xs: 'w-3 h-3',
  sm: 'w-4 h-4',
  md: 'w-6 h-6',
  lg: 'w-8 h-8',
};

const borderClasses = {
  sm: 'border',
  md: 'border-2',
};

function Spinner({ size = 'sm', border = 'sm', className = '', colorClassName = 'border-geek-400' }) {
  const resolvedSize = sizeClasses[size] || sizeClasses.sm;
  const resolvedBorder = borderClasses[border] || borderClasses.sm;

  return (
    <div
      className={[
        resolvedSize,
        resolvedBorder,
        colorClassName,
        'border-t-transparent rounded-full animate-spin',
        className,
      ].join(' ')}
    />
  );
}

export default Spinner;
