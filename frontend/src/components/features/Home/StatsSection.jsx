import { motion } from 'motion/react';

function StatsSection({ stats, isLoading }) {
  if (!isLoading && (!stats || stats.length === 0)) {
    return null;
  }

  return (
    <div className="py-10 px-4 md:px-8 border-y border-neutral-700/50">
      <div className="w-full max-w-[1200px] mx-auto">
        <div className="grid grid-cols-2 md:grid-cols-4 divide-x divide-neutral-700/50">
          {isLoading
            ? Array.from({ length: 4 }).map((_, index) => (
                <div key={index} className="px-6 py-4 flex flex-col gap-2">
                  <div className="h-8 w-20 rounded bg-neutral-700/60 animate-pulse" />
                  <div className="h-4 w-24 rounded bg-neutral-800/60 animate-pulse" />
                </div>
              ))
            : stats.map((stat, index) => (
                <motion.div
                  key={index}
                  className="px-6 py-4"
                  initial={{ opacity: 0 }}
                  whileInView={{ opacity: 1 }}
                  viewport={{ once: true }}
                  transition={{ delay: index * 0.06, duration: 0.35, ease: [0.16, 1, 0.3, 1] }}
                >
                  <span className="block text-3xl font-mono text-geek-400 leading-none mb-1">{stat.value}</span>
                  <span className="block text-sm text-neutral-400 font-mono tracking-wide">{stat.label}</span>
                </motion.div>
              ))}
        </div>
      </div>
    </div>
  );
}

export default StatsSection;
