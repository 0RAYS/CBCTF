import { motion } from 'motion/react';
import { ScrollingText, Button } from '../../../components/common';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useBranding } from '../../../hooks/useBranding';

const RANK_COLORS = ['text-rank-gold', 'text-rank-silver', 'text-rank-bronze'];

function LeaderboardPreview({ topUsers, isLoading }) {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { home } = useBranding();

  return (
    <div className="py-16 md:py-24 px-4 md:px-8">
      <div className="w-full max-w-[1200px] mx-auto">
        {/* Left-aligned heading — same rhythm as UpcomingContests */}
        <motion.div
          className="mb-8"
          initial={{ opacity: 0, y: 16 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ ease: [0.16, 1, 0.3, 1], duration: 0.4 }}
        >
          <div className="flex items-center gap-4 mb-3">
            <div className="w-6 h-[2px] bg-geek-400" />
            <span className="text-xs font-mono text-geek-400 tracking-[0.2em] uppercase">
              {home.leaderboard.titlePrefix || 'Rankings'}
            </span>
          </div>
          <h2 className="text-2xl md:text-3xl font-mono text-neutral-50">{home.leaderboard.titleHighlight}</h2>
          <p className="text-neutral-400 text-sm mt-2">{home.leaderboard.subtitle}</p>
        </motion.div>

        {/* Leaderboard rows — header row for structure */}
        <motion.div
          className="border border-neutral-600/60 rounded-md overflow-hidden"
          initial={{ opacity: 0, y: 16 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ ease: [0.16, 1, 0.3, 1], duration: 0.4 }}
        >
          {/* Column header */}
          <div className="grid grid-cols-[48px_1fr_auto_auto] gap-4 items-center px-5 py-2 border-b border-neutral-700/60 bg-neutral-800/40">
            <span className="text-[10px] font-mono text-neutral-500 uppercase tracking-[0.2em]">
              {t('common.scoreboard.headers.rank')}
            </span>
            <span className="text-[10px] font-mono text-neutral-500 uppercase tracking-[0.2em]">
              {t('common.scoreboard.headers.team')}
            </span>
            <span className="text-[10px] font-mono text-neutral-500 uppercase tracking-[0.2em] text-right">
              {t('common.scoreboard.headers.score')}
            </span>
            <span className="text-[10px] font-mono text-neutral-500 uppercase tracking-[0.2em] text-right w-16">
              {t('common.solved')}
            </span>
          </div>

          {isLoading
            ? Array.from({ length: 5 }).map((_, index) => (
                <div
                  key={index}
                  className={`h-[52px] bg-neutral-800/30 animate-pulse ${index !== 4 ? 'border-b border-neutral-700/50' : ''}`}
                />
              ))
            : topUsers?.length > 0 &&
              topUsers.map((team, index) => (
                <motion.div
                  key={index}
                  className={`grid grid-cols-[48px_1fr_auto_auto] gap-4 items-center px-5 py-3
                              ${index !== topUsers.length - 1 ? 'border-b border-neutral-700/50' : ''}
                              hover:bg-neutral-700/20 transition-colors duration-150 cursor-pointer group`}
                  initial={{ opacity: 0 }}
                  whileInView={{ opacity: 1 }}
                  viewport={{ once: true }}
                  transition={{ delay: index * 0.06, ease: [0.16, 1, 0.3, 1], duration: 0.3 }}
                >
                  <span className={`text-sm font-mono ${RANK_COLORS[index] ?? 'text-neutral-500'}`}>{index + 1}</span>
                  <ScrollingText
                    text={team.name}
                    className="text-sm text-neutral-200 font-mono group-hover:text-geek-400 transition-colors duration-150"
                    maxWidth={240}
                    speed={15}
                  />
                  <span className="text-sm font-mono text-geek-400 text-right">{team.score}</span>
                  <span className="text-sm font-mono text-neutral-300 text-right w-16">{team.solved}</span>
                </motion.div>
              ))}
        </motion.div>

        {!isLoading && (
          <div className="mt-6">
            <Button variant="outline" size="sm" onClick={() => navigate('/games')}>
              {home.leaderboard.action}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}

export default LeaderboardPreview;
