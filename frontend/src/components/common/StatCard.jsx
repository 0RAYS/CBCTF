import { motion } from 'motion/react';
import { EASE_T2 } from '../../config/motion';

function StatCard({ title, value, valueColor = 'text-geek-400', icon, iconBgClass = 'bg-geek-400/15', delay = 0 }) {
  return (
    <motion.div
      className="border border-neutral-600/50 rounded-md bg-neutral-800/50 p-4"
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      whileHover={{ y: -1, borderColor: 'rgba(89,126,247,0.35)' }}
      transition={{ duration: 0.22, delay, ease: EASE_T2 }}
    >
      {icon ? (
        <div className="flex items-center gap-3">
          <div className={`w-9 h-9 ${iconBgClass} rounded-md flex items-center justify-center flex-shrink-0`}>
            {icon}
          </div>
          <div>
            <p className="text-xs font-mono text-neutral-500 uppercase tracking-[0.12em] mb-0.5">{title}</p>
            <p className={`text-xl font-mono tabular-nums ${valueColor}`}>{value}</p>
          </div>
        </div>
      ) : (
        <>
          <p className="text-xs font-mono text-neutral-500 uppercase tracking-[0.12em] mb-1.5">{title}</p>
          <p className={`text-2xl font-mono tabular-nums ${valueColor}`}>{value}</p>
        </>
      )}
    </motion.div>
  );
}

export default StatCard;
