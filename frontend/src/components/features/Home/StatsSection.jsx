import { motion } from 'motion/react';
import { Card } from '../../../components/common';
import { useTranslation } from 'react-i18next';

function StatsSection(stats) {
  const { t } = useTranslation();
  stats = stats?.stats;
  if (stats?.length === 0) {
    stats = [
      { label: t('home.stats.activePlayers'), value: '10,000+' },
      { label: t('home.stats.challenges'), value: '500+' },
      { label: t('home.stats.events'), value: '50+' },
      { label: t('home.stats.successRate'), value: '85%' },
    ];
  }

  return (
    <div className="py-20 px-8">
      <div className="w-full max-w-[1200px] mx-auto">
        <motion.div
          className="grid grid-cols-4 gap-6"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.5 }}
        >
          {stats.length > 0 &&
            stats.map((stat, index) => (
              <Card
                key={index}
                variant="default"
                padding="md"
                className="flex flex-col items-center justify-center text-center hover:-translate-y-1 transition-transform duration-200"
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
