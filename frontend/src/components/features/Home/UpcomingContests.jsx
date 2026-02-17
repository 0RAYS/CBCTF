import { motion } from 'motion/react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

function UpcomingContests(contests) {
  // const contests = [
  //   {
  //     title: 'Web Warriors Challenge',
  //     date: '2024-04-15',
  //     duration: '48h',
  //     difficulty: 'Intermediate',
  //     prizes: '$5,000',
  //     registrations: 256,
  //     image: 'url_to_image',
  //   },
  //   {
  //     title: 'Crypto Masters Cup',
  //     date: '2024-04-20',
  //     duration: '24h',
  //     difficulty: 'Advanced',
  //     prizes: '$3,000',
  //     registrations: 128,
  //     image: 'url_to_image',
  //   },
  //   {
  //     title: 'Binary Blast',
  //     date: '2024-04-25',
  //     duration: '36h',
  //     difficulty: 'Expert',
  //     prizes: '$8,000',
  //     registrations: 192,
  //     image: 'url_to_image',
  //   },
  // ];
  const navigate = useNavigate();
  const { t } = useTranslation();
  contests = contests?.contests;

  return (
    <div className="py-20 px-8">
      <div className="w-full max-w-[1200px] mx-auto">
        {/* 标题 */}
        <motion.div
          className="text-center mb-12"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
        >
          <h2 className="text-3xl font-mono text-neutral-50 mb-4">
            {t('home.upcoming.titlePrefix')} <span className="text-geek-400">{t('home.upcoming.titleHighlight')}</span>
          </h2>
          <p className="text-neutral-300">{t('home.upcoming.subtitle')}</p>
        </motion.div>

        {/* 比赛列表 */}
        <div className={`grid grid-cols-${contests.length} gap-6`}>
          {contests.length > 0 &&
            contests.map((contest, index) => (
              <motion.div
                key={index}
                className="border border-neutral-300 rounded-md overflow-hidden
                                bg-black/30 backdrop-blur-[2px] group cursor-pointer"
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
                whileHover={{ y: -5 }}
                onClick={() => {
                  navigate(`/games`);
                }}
              >
                {/* 比赛信息 */}
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
      </div>
    </div>
  );
}

export default UpcomingContests;
