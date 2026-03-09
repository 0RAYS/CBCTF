import { motion } from 'motion/react';
import { ScrollingText, Button } from '../../../components/common';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

function LeaderboardPreview({ topUsers, isLoading }) {
  const navigate = useNavigate();
  const { t } = useTranslation();

  return (
    <div className="py-12 md:py-20 px-4 md:px-8">
      <div className="w-full max-w-[1200px] mx-auto">
        {/* Leaderboard */}
        <motion.div
          className="border border-neutral-300/30 rounded-md overflow-hidden"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
        >
          {isLoading
            ? Array.from({ length: 5 }).map((_, index) => (
                <div
                  key={index}
                  className={`h-[72px] bg-neutral-800 animate-pulse ${index !== 4 ? 'border-b border-neutral-600' : ''}`}
                />
              ))
            : topUsers?.length > 0 &&
              topUsers.map((team, index) => (
                <motion.div
                  key={index}
                  className={`flex items-center justify-between p-6 bg-neutral-900
                              ${index !== topUsers.length - 1 ? 'border-b border-neutral-600' : ''}
                              hover:bg-neutral-800 transition-colors duration-200 cursor-pointer group`}
                  initial={{ opacity: 0, x: -20 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: index * 0.1 }}
                >
                  {/* Rank and team info */}
                  <div className="flex items-center gap-6">
                    <span
                      className={`text-2xl font-mono ${
                        index === 0
                          ? 'text-yellow-400'
                          : index === 1
                            ? 'text-neutral-300'
                            : index === 2
                              ? 'text-yellow-600'
                              : 'text-neutral-400'
                      }`}
                    >
                      #{index + 1}
                    </span>
                    <div className="flex items-center gap-2">
                      <ScrollingText
                        text={team.name}
                        className="text-neutral-50 font-mono group-hover:text-geek-400 transition-colors duration-200"
                        maxWidth={200}
                        speed={15}
                      />
                    </div>
                  </div>

                  {/* Score and solved count */}
                  <div className="flex items-center gap-8">
                    <div className="flex items-center gap-2">
                      <span className="text-neutral-400">{t('common.score')}:</span>
                      <span className="text-geek-400 font-mono">{team.score}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-neutral-400">{t('common.solved')}:</span>
                      <span className="text-neutral-50 font-mono">{team.solved}</span>
                    </div>
                  </div>
                </motion.div>
              ))}
        </motion.div>

        {/* CTA */}
        {!isLoading && (
          <div className="flex justify-center mt-8">
            <Button variant="outline" onClick={() => navigate('/games')}>
              {t('home.leaderboard.viewScoreboard')}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}

export default LeaderboardPreview;
