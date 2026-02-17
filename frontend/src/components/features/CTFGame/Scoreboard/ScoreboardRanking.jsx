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
          className="grid gap-4 p-4 border-b border-neutral-300/30 place-items-center"
          style={{ gridTemplateColumns: gridCols }}
        >
          <div className="text-neutral-400 font-mono text-sm">{resolvedLabels.rank}</div>
          <div className="text-neutral-400 font-mono text-sm">{resolvedLabels.team}</div>
          <div className="text-neutral-400 font-mono text-sm flex items-center justify-end">{resolvedLabels.score}</div>
          <div className="text-neutral-400 font-mono text-sm flex items-center justify-center">
            {resolvedLabels.challenges}
          </div>
          <div className="text-neutral-400 font-mono text-sm flex items-center justify-end">
            {resolvedLabels.lastSubmit}
          </div>
        </div>

        <div className="divide-y divide-neutral-300/10">
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
                  className={`grid gap-4 p-4 hover:bg-neutral-300/5 transition-colors duration-200 ${onRowClick ? 'cursor-pointer' : ''}`}
                  style={{ gridTemplateColumns: gridCols }}
                  onClick={() => onRowClick && onRowClick(team, index)}
                >
                  <div className="flex items-center justify-center">
                    <span
                      className={`font-mono text-xl ${
                        rankValue === 1
                          ? 'text-yellow-400'
                          : rankValue === 2
                            ? 'text-neutral-300'
                            : rankValue === 3
                              ? 'text-orange-400'
                              : 'text-neutral-400'
                      }`}
                    >
                      #{rankValue}
                    </span>
                  </div>

                  <div className="flex items-center gap-4 justify-start">
                    <Avatar src={team.picture} name={team.name} size="xs" className="border border-neutral-300/30" />
                    <ScrollingText text={team.name} className="text-neutral-50 font-mono" maxWidth={240} speed={15} />
                  </div>

                  <div className="flex items-center justify-center">
                    <span className="text-geek-400 font-mono">{scoreValue}</span>
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
