import { motion } from 'motion/react';
import { Card, EmptyState, Avatar } from '../../common';
import ChallengeSolves from './ChallengeSolves';

const DEFAULT_LABELS = {
  rank: 'RANK',
  team: 'TEAM',
  score: 'SCORE',
  challenges: 'CHALLENGES',
  lastSubmit: 'LAST SUBMIT',
  total: 'TOTAL',
};

function ScoreboardRanking({ teams = [], labels = {}, locale = 'en-US', emptyMessage, footer = null, onRowClick }) {
  const resolvedLabels = { ...DEFAULT_LABELS, ...labels };
  const resolvedEmptyMessage = emptyMessage || 'No data';
  const gridCols =
    'grid-cols-[2rem_minmax(0,1fr)_minmax(0,0.65fr)] lg:grid-cols-[3rem_minmax(0,1.4fr)_minmax(0,0.8fr)_minmax(0,2fr)_minmax(0,1.2fr)]';

  return (
    <div className="w-full min-w-0 space-y-6">
      <Card variant="default" padding="none" animate className="overflow-hidden">
        <div
          className={`grid ${gridCols} gap-2 lg:gap-4 p-3 border-b border-neutral-600/50 place-items-center bg-neutral-800/40`}
        >
          <div className="text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase">
            {resolvedLabels.rank}
          </div>
          <div className="text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase">
            {resolvedLabels.team}
          </div>
          <div className="text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase flex items-center justify-end">
            {resolvedLabels.score}
          </div>
          <div className="hidden lg:flex text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase items-center justify-center">
            {resolvedLabels.challenges}
          </div>
          <div className="hidden lg:flex text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase items-center justify-end">
            {resolvedLabels.lastSubmit}
          </div>
        </div>

        <div className="overflow-hidden">
          {teams.length === 0 ? (
            <div className="py-10">
              <EmptyState title={resolvedEmptyMessage} />
            </div>
          ) : (
            teams.map((team, index) => {
              const rankValue = team.rank ?? index + 1;
              const scoreValue = typeof team.score === 'number' ? team.score.toLocaleString(locale) : team.score;
              const hasSolved = Array.isArray(team.solved) && team.solved.length > 0;
              const lastSubmit = team.lastSubmit || '-';

              return (
                <motion.div
                  key={team.id || team.name || index}
                  className={`grid ${gridCols} gap-2 lg:gap-4 p-3 transition-colors duration-200 ${onRowClick ? 'cursor-pointer focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-geek-400' : ''}
                    ${
                      rankValue === 1
                        ? 'bg-rank-gold/5 hover:bg-rank-gold/8 border-b border-rank-gold/15'
                        : rankValue === 2
                          ? 'bg-rank-silver/4 hover:bg-rank-silver/7 border-b border-neutral-600/30'
                          : rankValue === 3
                            ? 'bg-rank-bronze/4 hover:bg-rank-bronze/7 border-b border-neutral-600/30'
                            : 'hover:bg-neutral-300/5 border-b border-neutral-700/30 last:border-b-0'
                    }`}
                  onClick={() => onRowClick && onRowClick(team, index)}
                  role={onRowClick ? 'button' : undefined}
                  tabIndex={onRowClick ? 0 : undefined}
                  onKeyDown={
                    onRowClick
                      ? (event) => {
                          if (event.key === 'Enter' || event.key === ' ') {
                            event.preventDefault();
                            onRowClick(team, index);
                          }
                        }
                      : undefined
                  }
                >
                  <div className="flex items-center justify-center">
                    <span
                      className={`font-mono tabular-nums leading-none break-all ${
                        rankValue === 1
                          ? 'text-rank-gold text-2xl font-bold'
                          : rankValue === 2
                            ? 'text-rank-silver text-xl font-semibold'
                            : rankValue === 3
                              ? 'text-rank-bronze text-lg font-semibold'
                              : 'text-neutral-500 text-sm'
                      }`}
                    >
                      {rankValue}
                    </span>
                  </div>

                  <div className="flex min-w-0 items-center gap-2 lg:gap-3 justify-start">
                    <Avatar
                      src={team.picture}
                      name={team.name}
                      size="xs"
                      className="shrink-0 border border-neutral-300/30"
                    />
                    <span className="min-w-0 text-sm text-neutral-50 font-mono [overflow-wrap:anywhere]">
                      {team.name}
                    </span>
                  </div>

                  <div className="flex min-w-0 items-center justify-end lg:justify-center">
                    <span className="text-geek-400 text-sm lg:text-base font-mono tabular-nums break-all">
                      {scoreValue}
                    </span>
                  </div>

                  <div className="col-span-3 min-w-0 lg:col-span-1">
                    <div className="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs font-mono text-neutral-400 lg:hidden">
                      <span>{resolvedLabels.challenges}:</span>
                      {hasSolved ? (
                        team.solved.map((category) => (
                          <span key={category.category} className="[overflow-wrap:anywhere]">
                            {category.category}: {category.solved}/{category.all}
                          </span>
                        ))
                      ) : (
                        <span>-</span>
                      )}
                      {hasSolved && (
                        <span className="text-geek-400">
                          {resolvedLabels.total}: {team.totalSolved}
                        </span>
                      )}
                    </div>
                    <div className="hidden lg:block [&>div]:flex-wrap">
                      {hasSolved ? (
                        <ChallengeSolves
                          solved={team.solved}
                          totalSolved={team.totalSolved}
                          totalLabel={resolvedLabels.total}
                        />
                      ) : (
                        <div className="flex items-center justify-center text-neutral-500 font-mono text-sm">-</div>
                      )}
                    </div>
                  </div>

                  <div className="col-span-3 min-w-0 flex flex-wrap items-center gap-x-2 lg:col-span-1 lg:justify-end">
                    <span className="text-neutral-500 font-mono text-xs lg:hidden">{resolvedLabels.lastSubmit}:</span>
                    <span className="text-neutral-400 font-mono text-xs lg:text-sm [overflow-wrap:anywhere]">
                      {lastSubmit}
                    </span>
                  </div>
                </motion.div>
              );
            })
          )}
        </div>
      </Card>

      {footer}
    </div>
  );
}

export default ScoreboardRanking;
