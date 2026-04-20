import { motion } from 'motion/react';
import { ScrollingText, Card, EmptyState, Avatar } from '../../../../components/common';
import ChallengeSolves from '../ChallengeSolves';

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
  const gridCols = '60px 240px 160px auto 180px';

  return (
    <div className="w-full space-y-6">
      <Card variant="default" padding="none" animate className="overflow-hidden">
        <div
          className="grid gap-4 p-3 border-b border-neutral-600/50 place-items-center bg-neutral-800/40"
          style={{ gridTemplateColumns: gridCols }}
        >
          <div className="text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase">{resolvedLabels.rank}</div>
          <div className="text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase">{resolvedLabels.team}</div>
          <div className="text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase flex items-center justify-end">{resolvedLabels.score}</div>
          <div className="text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase flex items-center justify-center">
            {resolvedLabels.challenges}
          </div>
          <div className="text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase flex items-center justify-end">
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
                  className={`grid gap-4 p-4 transition-colors duration-200 ${onRowClick ? 'cursor-pointer' : ''}
                    ${
                      rankValue === 1
                        ? 'bg-rank-gold/5 hover:bg-rank-gold/8 border-b border-rank-gold/15'
                        : rankValue === 2
                          ? 'bg-rank-silver/4 hover:bg-rank-silver/7 border-b border-neutral-600/30'
                          : rankValue === 3
                            ? 'bg-rank-bronze/4 hover:bg-rank-bronze/7 border-b border-neutral-600/30'
                            : 'hover:bg-neutral-300/5 border-b border-neutral-700/30 last:border-b-0'
                    }`}
                  style={{ gridTemplateColumns: gridCols }}
                  onClick={() => onRowClick && onRowClick(team, index)}
                >
                  <div className="flex items-center justify-center">
                    <span
                      className={`font-mono tabular-nums leading-none ${
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

                  <div className="flex items-center gap-4 justify-start">
                    <Avatar src={team.picture} name={team.name} size="xs" className="border border-neutral-300/30" />
                    <ScrollingText text={team.name} className="text-neutral-50 font-mono" maxWidth={240} speed={15} />
                  </div>

                  <div className="flex items-center justify-center">
                    <span className="text-geek-400 font-mono tabular-nums">{scoreValue}</span>
                  </div>

                  {hasSolved ? (
                    <ChallengeSolves
                      solved={team.solved}
                      totalSolved={team.totalSolved}
                      totalLabel={resolvedLabels.total}
                    />
                  ) : (
                    <div className="flex items-center justify-center text-neutral-500 font-mono text-sm">-</div>
                  )}

                  <div className="flex items-center justify-end">
                    <span className="text-neutral-400 font-mono text-sm">{lastSubmit}</span>
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
