import { motion } from 'motion/react';
import { Card } from '../../../components/common';
import { useTranslation } from 'react-i18next';

function StatsSection({ stats, isLoading }) {
  const { t } = useTranslation();

  if (!isLoading && (!stats || stats.length === 0)) {
    return null;
  }

  return (
    <div className="py-12 md:py-20 px-4 md:px-8">
      <div className="w-full max-w-[1200px] mx-auto">
        <motion.div
          className="grid grid-cols-2 md:grid-cols-4 gap-6"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.5 }}
        >
          {isLoading
            ? Array.from({ length: 4 }).map((_, index) => (
                <div
                  key={index}
                  className="h-[88px] rounded-md bg-neutral-800 animate-pulse"
                />
              ))
            : stats.map((stat, index) => (
                <Card
                  key={index}
                  variant="default"
                  padding="md"
                  className="flex flex-col items-center justify-center text-center hover:border-geek-400/50 transition-colors duration-200"
                >
                  <span className="text-3xl font-mono text-geek-400 mb-2">{stat.value}</span>
                  <span className="text-neutral-300">{stat.label}</span>
                </Card>
              ))}
        </motion.div>
      </div>
    </div>
  );
}

export default StatsSection;
