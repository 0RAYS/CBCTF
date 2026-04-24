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

/**
 * ChallengeSkeleton — ghost card matching ChallengeBoard card layout
 */
export function ChallengeSkeleton() {
  return (
    <div className="p-4 border border-neutral-700/40 rounded-md bg-neutral-800/30">
      {/* Title row */}
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-3 flex-1">
          <Skeleton variant="line" className="w-14 h-3" />
          <Skeleton variant="line" className="w-28 h-3" />
          <Skeleton variant="line" className="w-10 h-3 ml-auto" />
        </div>
        <div className="flex items-center gap-2 ml-4">
          <Skeleton variant="line" className="w-12 h-3" />
        </div>
      </div>
      {/* Tags row */}
      <div className="flex items-center gap-2">
        <Skeleton variant="line" className="w-12 h-4" />
        <Skeleton variant="line" className="w-16 h-4" />
      </div>
    </div>
  );
}

export default Skeleton;
