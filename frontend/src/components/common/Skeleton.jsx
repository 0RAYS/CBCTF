/**
 * Skeleton — loading placeholder with animated pulse.
 * Use to replace content while data is loading — prevents jarring blank-flash.
 *
 * @param {string} className - Tailwind classes for sizing/shape
 * @param {'line'|'rect'|'circle'} variant - Semantic shape shorthand
 */
function Skeleton({ className = '', variant = 'rect' }) {
  const base = 'animate-pulse bg-neutral-700/40 rounded';

  const variantClass = variant === 'line' ? 'h-3 rounded-full' : variant === 'circle' ? 'rounded-full' : 'rounded-md';

  return <div className={`${base} ${variantClass} ${className}`} aria-hidden="true" />;
}

export default Skeleton;
