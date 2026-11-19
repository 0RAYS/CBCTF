/** Chip — inline label badge; variant="pill" 是实底圆角, variant="tag" 是描边, colorClass 覆盖默认配色 */
function Chip({ label, variant = 'pill', size = 'md', colorClass, className = '', title }) {
  const py = size === 'sm' ? 'py-0.5' : 'py-1';

  if (variant === 'tag') {
    return (
      <span
        className={`px-2 ${py} bg-neutral-700/40 border rounded text-xs font-mono ${colorClass || 'border-neutral-600/60 text-neutral-300'} ${className}`}
        title={title}
      >
        {label}
      </span>
    );
  }

  return (
    <span
      className={`px-2 ${py} rounded-full text-xs font-mono ${colorClass || 'bg-neutral-400/20 text-neutral-400'} ${className}`}
      title={title}
    >
      {label}
    </span>
  );
}

export default Chip;
