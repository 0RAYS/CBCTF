import { motion } from 'motion/react';

/**
 * StatCard — bordered metric card
 *
 * Without icon: title label above a large value
 * With icon:    icon thumbnail left, title + value right
 *
 * Props:
 *   title       {string}       — label / description
 *   value       {string|number}— metric value
 *   valueColor  {string}       — Tailwind text colour class (default 'text-geek-400')
 *   icon        {ReactNode}    — optional icon element
 *   iconBgClass {string}       — icon container bg class (default 'bg-geek-400/20')
 *   delay       {number}       — motion entrance delay in seconds (default 0)
 */
function StatCard({ title, value, valueColor = 'text-geek-400', icon, iconBgClass = 'bg-geek-400/20', delay = 0 }) {
  return (
    <motion.div
      className="border border-neutral-600 rounded-md bg-neutral-900 p-4"
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      whileHover={{ y: -2, boxShadow: '0 4px 20px rgba(89, 126, 247, 0.1)' }}
      transition={{ duration: 0.2, delay }}
    >
      {icon ? (
        <div className="flex items-center gap-3">
          <div className={`w-10 h-10 ${iconBgClass} rounded-md flex items-center justify-center flex-shrink-0`}>
            {icon}
          </div>
          <div>
            <p className="text-sm font-mono text-neutral-400">{title}</p>
            <p className={`text-2xl font-mono ${valueColor}`}>{value}</p>
          </div>
        </div>
      ) : (
        <>
          <h2 className="text-sm font-mono text-neutral-400 mb-2">{title}</h2>
          <p className={`text-2xl font-mono ${valueColor}`}>{value}</p>
        </>
      )}
    </motion.div>
  );
}

export default StatCard;
