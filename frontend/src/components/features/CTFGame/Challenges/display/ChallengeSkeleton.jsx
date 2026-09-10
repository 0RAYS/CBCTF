import Skeleton from '../../../../common/Skeleton';

export default function ChallengeSkeleton() {
  return (
    <div className="p-4 border border-neutral-700/40 rounded-md bg-neutral-800/30">
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
      <div className="flex items-center gap-2">
        <Skeleton variant="line" className="w-12 h-4" />
        <Skeleton variant="line" className="w-16 h-4" />
      </div>
    </div>
  );
}
