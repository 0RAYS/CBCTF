import { motion } from 'motion/react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { IconCheck } from '@tabler/icons-react';
import { Avatar, Card, Pagination, ScrollingText } from '@/components/common';
import Loading from '@/components/common/Loading';
import { useTranslation } from 'react-i18next';
import { buildLeaderboardColumns, indexTeamChallenges } from '@/components/features/Leaderboard/leaderboardModel.js';
import useScoreColumnWidth from '@/components/features/Leaderboard/useScoreColumnWidth';
import TruncatedText from '@/components/common/TruncatedText';

/** 通用表格式积分榜组件 */
function LeaderboardTable({
  teams = [],
  challenges = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 20,
  onPageChange,
  loading = false,
}) {
  const { t } = useTranslation();
  const containerRef = useRef(null);
  const [stickyWidths, setStickyWidths] = useState({
    containerWidth: 0,
    rankWidth: 60,
    teamWidth: 240,
  });
  const scoreValues = useMemo(() => teams.map((team) => team.score), [teams]);
  const { width: scoreWidth, measurement } = useScoreColumnWidth({
    scores: scoreValues,
    label: t('common.leaderboard.headers.score'),
    scoreClassName: 'text-base',
    padding: 25,
    minWidth: 80,
  });

  useEffect(() => {
    const container = containerRef.current;
    if (!container || typeof ResizeObserver === 'undefined') {
      return;
    }

    const updateWidths = () => {
      const width = container.clientWidth || 0;
      if (!width) return;
      const rankWidth = width < 640 ? 52 : 60;
      const teamWidth = Math.max(140, Math.min(220, Math.floor(width * 0.24)));

      setStickyWidths((prev) => {
        if (prev.rankWidth === rankWidth && prev.teamWidth === teamWidth && prev.containerWidth === width) {
          return prev;
        }
        return { containerWidth: width, rankWidth, teamWidth };
      });
    };

    updateWidths();
    const observer = new ResizeObserver(updateWidths);
    observer.observe(container);

    return () => {
      observer.disconnect();
    };
  }, []);

  const {
    categoryEntries,
    columns: challengeColumns,
    widths: challengeColWidths,
    widthById: challengeWidthMap,
  } = useMemo(() => buildLeaderboardColumns(challenges), [challenges]);
  const teamChallengeIndexes = useMemo(() => teams.map((team) => indexTeamChallenges(team.challenges || [])), [teams]);
  const challengeWidth = challengeColWidths.reduce((total, width) => total + width, 0);
  // Fill spare space with the team column so actual widths keep matching sticky offsets.
  const teamWidth = Math.max(
    stickyWidths.teamWidth,
    stickyWidths.containerWidth - stickyWidths.rankWidth - scoreWidth - challengeWidth
  );
  const tableWidth = stickyWidths.rankWidth + teamWidth + scoreWidth + challengeWidth;
  // On narrow screens, let team/score scroll so pinned columns cannot cover every challenge.
  const detailPosition =
    stickyWidths.rankWidth + teamWidth + scoreWidth + 80 <= stickyWidths.containerWidth ? 'sticky' : 'static';

  if (loading) {
    return (
      <div className="w-full min-w-0 space-y-6">
        {measurement}
        <Card variant="default" padding="none" className="overflow-hidden">
          <Loading compact />
        </Card>
      </div>
    );
  }

  return (
    <div className="w-full min-w-0 space-y-6">
      {measurement}
      <Card
        variant="default"
        padding="none"
        className="overflow-hidden"
        style={{
          '--lb-rank-width': `${stickyWidths.rankWidth}px`,
          '--lb-team-width': `${teamWidth}px`,
          '--lb-score-width': `${scoreWidth}px`,
        }}
      >
        <div className="overflow-x-auto" ref={containerRef}>
          <table style={{ width: tableWidth, tableLayout: 'fixed' }}>
            <colgroup>
              <col style={{ width: 'var(--lb-rank-width)' }} />
              <col style={{ width: 'var(--lb-team-width)' }} />
              <col style={{ width: 'var(--lb-score-width)' }} />
              {challengeColWidths.map((width, index) => (
                <col key={`challenge-col-${index}`} style={{ width: `${width}px` }} />
              ))}
            </colgroup>
            <thead>
              <tr className="bg-black/40">
                <th
                  scope="col"
                  className="sticky left-0 z-10 bg-neutral-800 p-3 text-center text-[10px] text-neutral-400 font-mono tracking-[0.18em] uppercase border-r border-neutral-600/40"
                  style={{ width: 'var(--lb-rank-width)', minWidth: 'var(--lb-rank-width)' }}
                >
                  <TruncatedText>{t('common.leaderboard.headers.rank')}</TruncatedText>
                </th>
                <th
                  scope="col"
                  className="z-10 bg-neutral-800 p-3 text-center text-[10px] text-neutral-400 font-mono tracking-[0.18em] uppercase border-r border-neutral-600/40"
                  style={{
                    position: detailPosition,
                    width: 'var(--lb-team-width)',
                    minWidth: 'var(--lb-team-width)',
                    left: 'var(--lb-rank-width)',
                  }}
                >
                  <TruncatedText>{t('common.leaderboard.headers.team')}</TruncatedText>
                </th>
                <th
                  scope="col"
                  className="z-10 bg-neutral-800 p-3 text-center text-[10px] text-neutral-400 font-mono tracking-[0.18em] uppercase whitespace-nowrap border-r border-neutral-600/40"
                  style={{
                    position: detailPosition,
                    width: 'var(--lb-score-width)',
                    minWidth: 'var(--lb-score-width)',
                    left: 'calc(var(--lb-rank-width) + var(--lb-team-width))',
                  }}
                >
                  <TruncatedText>{t('common.leaderboard.headers.score')}</TruncatedText>
                </th>

                {categoryEntries.map(([category, categoryIchallenges]) => (
                  <th
                    key={category}
                    scope="colgroup"
                    className="text-center text-neutral-400 font-mono text-sm border-r border-neutral-300/30"
                    colSpan={categoryIchallenges.length}
                  >
                    <div className="flex items-center justify-center gap-2 p-3">
                      <span className="shrink-0 text-geek-400">#</span>
                      <TruncatedText maxWidth={240}>{category}</TruncatedText>
                      <span className="shrink-0 text-xs text-neutral-500">({categoryIchallenges.length})</span>
                    </div>
                  </th>
                ))}
              </tr>

              <tr className="bg-black/30 border-t border-neutral-600/40">
                <th
                  className="sticky left-0 z-10 bg-neutral-800 p-2 border-r border-neutral-600/40"
                  style={{ width: 'var(--lb-rank-width)', minWidth: 'var(--lb-rank-width)' }}
                ></th>
                <th
                  className="z-10 bg-neutral-800 p-2 border-r border-neutral-600/40"
                  style={{
                    position: detailPosition,
                    width: 'var(--lb-team-width)',
                    minWidth: 'var(--lb-team-width)',
                    left: 'var(--lb-rank-width)',
                  }}
                ></th>
                <th
                  className="z-10 bg-neutral-800 p-2 border-r border-neutral-600/40"
                  style={{
                    position: detailPosition,
                    width: 'var(--lb-score-width)',
                    minWidth: 'var(--lb-score-width)',
                    left: 'calc(var(--lb-rank-width) + var(--lb-team-width))',
                  }}
                ></th>

                {categoryEntries.map(([, categoryIchallenges]) =>
                  categoryIchallenges.map((challenge) => {
                    const colWidth = challengeWidthMap.get(challenge.id) || 80;
                    return (
                      <th
                        scope="col"
                        key={challenge.id}
                        className="p-3 text-center text-neutral-400 font-mono text-xs border-r border-neutral-300/30"
                        style={{ width: `${colWidth}px` }}
                        title={challenge.name}
                      >
                        <TruncatedText>{challenge.name}</TruncatedText>
                      </th>
                    );
                  })
                )}
              </tr>
            </thead>

            <tbody>
              {teams.map((team, teamIndex) => (
                <motion.tr
                  key={team.id}
                  className="border-t border-neutral-300/10 hover:bg-black/40"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: teamIndex * 0.05 }}
                >
                  <td
                    className="sticky left-0 z-10 bg-neutral-800 p-3 text-center text-neutral-300 font-mono tabular-nums border-r border-neutral-600/40"
                    style={{ width: 'var(--lb-rank-width)', minWidth: 'var(--lb-rank-width)' }}
                  >
                    <TruncatedText>{team.rank ?? (currentPage - 1) * pageSize + teamIndex + 1}</TruncatedText>
                  </td>

                  <td
                    className="z-10 bg-neutral-800 p-3 border-r border-neutral-600/40"
                    style={{
                      position: detailPosition,
                      width: 'var(--lb-team-width)',
                      minWidth: 'var(--lb-team-width)',
                      left: 'var(--lb-rank-width)',
                    }}
                  >
                    <div className="flex items-center gap-2">
                      <Avatar src={team.picture} name={team.name} size="xs" className="border border-neutral-600/60" />
                      <div className="min-w-0 flex-1">
                        <ScrollingText
                          text={team.name}
                          className="text-neutral-50 font-mono text-sm"
                          maxWidth={240}
                          speed={15}
                        />
                        <TruncatedText className="text-xs text-neutral-500 tabular-nums">
                          {t('game.leaderboardTable.members', { count: team.users })}
                        </TruncatedText>
                      </div>
                    </div>
                  </td>

                  <td
                    className="z-10 bg-neutral-800 p-3 text-base text-center text-geek-400 font-mono tabular-nums whitespace-nowrap border-r border-neutral-600/40"
                    style={{
                      position: detailPosition,
                      width: 'var(--lb-score-width)',
                      minWidth: 'var(--lb-score-width)',
                      left: 'calc(var(--lb-rank-width) + var(--lb-team-width))',
                    }}
                  >
                    <TruncatedText>{scoreValues[teamIndex]}</TruncatedText>
                  </td>

                  {challengeColumns.map((challenge) => {
                    const solved = teamChallengeIndexes[teamIndex].get(challenge.id)?.solved > 0;

                    return (
                      <td key={challenge.id} className="group/cell p-3 text-center border-r border-neutral-300/30">
                        <div
                          className="transition-transform duration-200 group-hover/cell:scale-110 flex items-center justify-center font-mono"
                          title={`${team.name} - ${challenge.name}: ${
                            solved
                              ? t('game.leaderboardTable.status.solved')
                              : t('game.leaderboardTable.status.unsolved')
                          }`}
                        >
                          {solved ? (
                            <IconCheck size={18} stroke={2} className="text-green-400" />
                          ) : (
                            <span className="text-neutral-500">•</span>
                          )}
                        </div>
                      </td>
                    );
                  })}
                </motion.tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>

      <div className="mt-6 flex justify-center w-full overflow-x-auto">
        <Pagination
          current={currentPage}
          total={Math.ceil(totalCount / pageSize)}
          onChange={onPageChange}
          showTotal={true}
          totalItems={totalCount}
        />
      </div>
    </div>
  );
}

export default LeaderboardTable;
