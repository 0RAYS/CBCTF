import { motion } from 'motion/react';
import { Card } from '../../../components/common';

function StatsSection({ stats, isLoading }) {
  if (!isLoading && (!stats || stats.length === 0)) {
    return null;
  }

  return (
    <div className="py-12 md:py-20 px-4 md:px-8">
      <div className="w-full max-w-[1200px] mx-auto">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
          {isLoading
            ? Array.from({ length: 4 }).map((_, index) => (
                <div key={index} className="h-[88px] rounded-md bg-neutral-800 animate-pulse" />
              ))
            : stats.map((stat, index) => (
                <motion.div
                  key={index}
                  initial={{ opacity: 0, y: 16 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: index * 0.08, duration: 0.4, ease: [0.25, 1, 0.5, 1] }}
                >
                  <Card
                    variant="default"
                    padding="md"
                    className="flex flex-col items-center justify-center text-center hover:border-geek-400/50 transition-colors duration-200"
                  >
                    <span className="text-3xl font-mono text-geek-400 mb-2">{stat.value}</span>
                    <span className="text-neutral-300">{stat.label}</span>
                  </Card>
                </motion.div>
              ))}
        </div>
      </div>
    </div>
  );
}

export default StatsSection;
