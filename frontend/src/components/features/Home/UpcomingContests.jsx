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
    <div className="py-12 md:py-20 px-4 md:px-8">
      <div className="w-full max-w-[1200px] mx-auto">
        {/* Heading */}
        <motion.div
          className="text-center mb-12"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ ease: [0.25, 1, 0.5, 1], duration: 0.4 }}
        >
          <h2 className="text-3xl font-mono text-neutral-50 mb-4">
            {home.upcoming.titlePrefix} <span className="text-geek-400">{home.upcoming.titleHighlight}</span>
          </h2>
          <p className="text-neutral-300">{home.upcoming.subtitle}</p>
        </motion.div>

        {/* Contest list */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {isLoading
            ? Array.from({ length: 3 }).map((_, index) => (
                <div key={index} className="h-[200px] rounded-md bg-neutral-800 animate-pulse" />
              ))
            : contests.map((contest, index) => (
                <motion.div
                  key={index}
                  className="border border-neutral-300/30 rounded-md overflow-hidden
                             bg-neutral-900 group cursor-pointer hover:border-geek-400/50 transition-colors duration-200"
                  initial={{ opacity: 0, y: 20 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: index * 0.1, ease: [0.25, 1, 0.5, 1], duration: 0.4 }}
                  onClick={() => navigate('/games')}
                >
                  {/* Contest info */}
                  <div className="p-6 space-y-4">
                    <h3
                      className="text-xl font-mono text-neutral-50 group-hover:text-geek-400
                                 transition-colors duration-200"
                    >
                      {contest.title}
                    </h3>

                    <div className="space-y-2">
                      <div className="flex items-center justify-between">
                        <span className="text-neutral-400">{t('common.date')}</span>
                        <span className="text-neutral-50 font-mono">{contest.date}</span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-neutral-400">{t('common.duration')}</span>
                        <span className="text-neutral-50 font-mono">{contest.duration}</span>
                      </div>
                    </div>

                    <div className="pt-4 border-t border-neutral-300/30 space-y-2">
                      <div className="flex items-center justify-between">
                        <span className="text-neutral-400">{t('common.registrations')}</span>
                        <span className="text-geek-400 font-mono">{contest.registrations}</span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-neutral-400">{t('common.teams')}</span>
                        <span className="text-geek-400 font-mono">{contest.teams}</span>
                      </div>
                    </div>
                  </div>
                </motion.div>
              ))}
        </div>

        {/* View All CTA */}
        {!isLoading && (
          <div className="flex justify-center mt-8">
            <Button variant="outline" onClick={() => navigate('/games')}>
              {home.upcoming.action}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}

export default UpcomingContests;
