/** Skeleton — loading placeholder with animated pulse */
function Skeleton({ className = '', variant = 'rect' }) {
  const base = 'animate-pulse bg-neutral-700/40 rounded';

  const variantClass = variant === 'line' ? 'h-3 rounded-full' : variant === 'circle' ? 'rounded-full' : 'rounded-md';

  return <div className={`${base} ${variantClass} ${className}`} aria-hidden="true" />;
}

export default Skeleton;
