import { motion } from 'motion/react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Button } from '../../../components/common';
import { useBranding } from '../../../hooks/useBranding';

function UpcomingContests({ contests = [], isLoading }) {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { home } = useBranding();

  return (
    <div className="py-16 md:py-24 px-4 md:px-8">
      <div className="w-full max-w-[1200px] mx-auto">
        {/* Left-aligned heading with accent rule */}
        <motion.div
          className="mb-10"
          initial={{ opacity: 0, y: 16 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ ease: [0.16, 1, 0.3, 1], duration: 0.4 }}
        >
          <div className="flex items-center gap-4 mb-3">
            <div className="w-6 h-[2px] bg-geek-400" />
            <span className="text-xs font-mono text-geek-400 tracking-[0.2em] uppercase">
              {home.upcoming.titlePrefix || 'Upcoming'}
            </span>
          </div>
          <h2 className="text-2xl md:text-3xl font-mono text-neutral-50">
            {home.upcoming.titleHighlight}
          </h2>
          <p className="text-neutral-400 text-sm mt-2 max-w-[60ch]">{home.upcoming.subtitle}</p>
        </motion.div>

        {/* Contest grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {isLoading
            ? Array.from({ length: 3 }).map((_, index) => (
                <div key={index} className="h-[180px] rounded-md bg-neutral-800/60 animate-pulse" />
              ))
            : contests.map((contest, index) => (
                <motion.div
                  key={index}
                  className="border border-neutral-600/60 rounded-md overflow-hidden
                             bg-neutral-800/40 group cursor-pointer hover:border-geek-400/50
                             hover:bg-neutral-800/60 transition-colors duration-200"
                  initial={{ opacity: 0, y: 16 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: index * 0.08, ease: [0.16, 1, 0.3, 1], duration: 0.35 }}
                  onClick={() => navigate('/games')}
                >
                  <div className="p-5 space-y-4">
                    <h3 className="text-base font-mono text-neutral-50 group-hover:text-geek-400 transition-colors duration-200 leading-snug">
                      {contest.title}
                    </h3>

                    <div className="space-y-1.5 text-sm font-mono">
                      <div className="flex items-baseline justify-between gap-2">
                        <span className="text-neutral-500">{t('common.date')}</span>
                        <span className="text-neutral-200">{contest.date}</span>
                      </div>
                      <div className="flex items-baseline justify-between gap-2">
                        <span className="text-neutral-500">{t('common.duration')}</span>
                        <span className="text-neutral-200">{contest.duration}</span>
                      </div>
                    </div>

                    <div className="pt-3 border-t border-neutral-600/50 flex items-center justify-between gap-4 text-sm font-mono">
                      <span className="text-neutral-500">{t('common.teams')}</span>
                      <span className="text-geek-400">{contest.teams}</span>
                    </div>
                  </div>
                </motion.div>
              ))}
        </div>

        {!isLoading && (
          <div className="mt-8">
            <Button variant="outline" size="sm" onClick={() => navigate('/games')}>
              {home.upcoming.action}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}

export default UpcomingContests;
