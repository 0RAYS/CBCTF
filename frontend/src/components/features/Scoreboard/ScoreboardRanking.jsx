import { motion } from 'motion/react';
import { useMemo } from 'react';
import { Avatar, Card, EmptyState } from '../../common';
import ChallengeSolves from './ChallengeSolves';
import useScoreColumnWidth from './useScoreColumnWidth';
import TruncatedText from '../../common/TruncatedText';

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
  const scoreValues = useMemo(
    () => teams.map((team) => (typeof team.score === 'number' ? team.score.toLocaleString(locale) : team.score)),
    [teams, locale]
  );
  const { width: scoreWidth, measurement } = useScoreColumnWidth({
    scores: scoreValues,
    label: resolvedLabels.score,
    scoreClassName: 'text-sm',
    padding: 4,
  });
  const gridCols =
    'grid-cols-[2rem_minmax(0,1fr)_var(--sb-score-width)] lg:grid-cols-[2.25rem_minmax(9rem,14rem)_var(--sb-score-width)_minmax(12rem,1fr)_9rem]';

  return (
    <div className="w-full min-w-0 space-y-4" style={{ '--sb-score-width': `${scoreWidth}px` }}>
      {measurement}
      <Card variant="default" padding="none" animate className="overflow-hidden">
        <div className="overflow-x-auto">
          <div className="min-w-[calc(var(--sb-score-width)+10rem)] lg:min-w-[calc(var(--sb-score-width)+37rem)]">
            <div
              className={`grid ${gridCols} gap-x-3 gap-y-1.5 px-3 py-2 border-b border-neutral-600/50 items-center bg-neutral-800/40`}
            >
              <div className="min-w-0 max-w-full text-center text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase">
                <TruncatedText>{resolvedLabels.rank}</TruncatedText>
              </div>
              <div className="min-w-0 max-w-full text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase">
                <TruncatedText>{resolvedLabels.team}</TruncatedText>
              </div>
              <div className="min-w-0 max-w-full text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase flex items-center justify-end">
                <TruncatedText>{resolvedLabels.score}</TruncatedText>
              </div>
              <div className="min-w-0 max-w-full hidden lg:flex text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase items-center">
                <TruncatedText>{resolvedLabels.challenges}</TruncatedText>
              </div>
              <div className="min-w-0 max-w-full hidden lg:flex text-[10px] font-mono text-neutral-500 tracking-[0.18em] uppercase items-center justify-end">
                <TruncatedText>{resolvedLabels.lastSubmit}</TruncatedText>
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
                  const scoreValue = scoreValues[index];
                  const hasSolved = Array.isArray(team.solved) && team.solved.length > 0;
                  const lastSubmit = team.lastSubmit || '-';

                  return (
                    <motion.div
                      key={team.id || team.name || index}
                      className={`grid ${gridCols} gap-x-3 gap-y-1.5 px-3 py-2 transition-colors duration-200 ${onRowClick ? 'cursor-pointer focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-geek-400' : ''}
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
                          title={String(rankValue)}
                          className={`min-w-0 max-w-full truncate font-mono tabular-nums leading-none ${
                            rankValue === 1
                              ? 'text-rank-gold text-lg font-bold'
                              : rankValue === 2
                                ? 'text-rank-silver text-base font-semibold'
                                : rankValue === 3
                                  ? 'text-rank-bronze text-base font-semibold'
                                  : 'text-neutral-500 text-sm'
                          }`}
                        >
                          {rankValue}
                        </span>
                      </div>

                      <div className="flex min-w-0 items-center gap-2 justify-start">
                        <Avatar
                          src={team.picture}
                          name={team.name}
                          size="xs"
                          className="shrink-0 border border-neutral-300/30"
                        />
                        <TruncatedText className="text-sm text-neutral-50 font-mono">{team.name}</TruncatedText>
                      </div>

                      <div className="flex min-w-0 items-center justify-end">
                        <TruncatedText className="text-geek-400 text-sm font-mono tabular-nums">
                          {scoreValue}
                        </TruncatedText>
                      </div>

                      <div
                        className="col-span-3 min-w-0 flex items-center lg:col-span-1"
                        role="group"
                        aria-label={resolvedLabels.challenges}
                      >
                        {hasSolved ? (
                          <ChallengeSolves
                            solved={team.solved}
                            totalSolved={team.totalSolved}
                            totalLabel={resolvedLabels.total}
                          />
                        ) : (
                          <div className="text-neutral-500 font-mono text-sm">-</div>
                        )}
                      </div>

                      <div className="col-span-3 min-w-0 flex flex-wrap items-center gap-x-2 lg:col-span-1 lg:justify-end">
                        <span className="text-neutral-500 font-mono text-xs lg:hidden">
                          {resolvedLabels.lastSubmit}:
                        </span>
                        <TruncatedText className="text-neutral-400 font-mono text-xs">{lastSubmit}</TruncatedText>
                      </div>
                    </motion.div>
                  );
                })
              )}
            </div>
          </div>
        </div>
      </Card>

      {footer}
    </div>
  );
}

export default ScoreboardRanking;
