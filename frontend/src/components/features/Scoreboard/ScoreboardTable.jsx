/**
 * 通用表格式积分榜组件
 * @param {Object} props
 * @param {Array} props.teams - 队伍列表
 * @param {Array} props.challenges - 题目列表（按分类分组）
 * @param {number} props.totalCount - 总队伍数
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {Function} props.onPageChange - 页码变化处理函数
 */

import { motion } from 'motion/react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { ScrollingText, Card, Avatar, Pagination } from '../../common';
import { useTranslation } from 'react-i18next';
import { buildScoreboardColumns, indexTeamChallenges } from './scoreboardModel.js';

function ScoreboardTable({
  teams = [],
  challenges = [],
  totalCount = 0,
  currentPage = 1,
  pageSize = 20,
  onPageChange,
}) {
  const { t } = useTranslation();
  const containerRef = useRef(null);
  const [stickyWidths, setStickyWidths] = useState({
    rankWidth: 60,
    teamWidth: 240,
    scoreWidth: 80,
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
      const scoreWidth = width < 640 ? 72 : 80;
      const teamWidth = Math.max(140, Math.min(220, Math.floor(width * 0.24)));

      setStickyWidths((prev) => {
        if (prev.rankWidth === rankWidth && prev.teamWidth === teamWidth && prev.scoreWidth === scoreWidth) {
          return prev;
        }
        return { rankWidth, teamWidth, scoreWidth };
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
  } = useMemo(() => buildScoreboardColumns(challenges), [challenges]);
  const teamChallengeIndexes = useMemo(() => teams.map((team) => indexTeamChallenges(team.challenges || [])), [teams]);

  // 计算排名
  const getTeamRank = (index) => {
    return (currentPage - 1) * pageSize + index + 1;
  };

  return (
    <div className="w-full space-y-6" ref={containerRef}>
      {/* 表格容器 */}
      <Card
        variant="default"
        padding="none"
        className="overflow-hidden"
        style={{
          '--sb-rank-width': `${stickyWidths.rankWidth}px`,
          '--sb-team-width': `${stickyWidths.teamWidth}px`,
          '--sb-score-width': `${stickyWidths.scoreWidth}px`,
        }}
      >
        <div className="overflow-x-auto">
          <table className="w-full" style={{ minWidth: '600px', tableLayout: 'fixed' }}>
            <colgroup>
              <col style={{ width: 'var(--sb-rank-width)' }} />
              <col style={{ width: 'var(--sb-team-width)' }} />
              <col style={{ width: 'var(--sb-score-width)' }} />
              {challengeColWidths.map((width, index) => (
                <col key={`challenge-col-${index}`} style={{ width: `${width}px` }} />
              ))}
            </colgroup>
            <thead>
              <tr className="bg-neutral-800/60">
                {/* 固定列 */}
                <th
                  scope="col"
                  className="sticky left-0 z-10 bg-neutral-800/80 p-3 text-center text-[10px] text-neutral-500 font-mono tracking-[0.18em] uppercase border-r border-neutral-600/40"
                  style={{ width: 'var(--sb-rank-width)', minWidth: 'var(--sb-rank-width)' }}
                >
                  {t('game.scoreboardTable.headers.rank')}
                </th>
                <th
                  scope="col"
                  className="sticky z-10 bg-neutral-800/80 p-3 text-center text-[10px] text-neutral-500 font-mono tracking-[0.18em] uppercase border-r border-neutral-600/40"
                  style={{
                    width: 'var(--sb-team-width)',
                    minWidth: 'var(--sb-team-width)',
                    left: 'var(--sb-rank-width)',
                  }}
                >
                  {t('game.scoreboardTable.headers.team')}
                </th>
                <th
                  scope="col"
                  className="sticky z-10 bg-neutral-800/80 p-3 text-center text-[10px] text-neutral-500 font-mono tracking-[0.18em] uppercase border-r border-neutral-600/40"
                  style={{
                    width: 'var(--sb-score-width)',
                    minWidth: 'var(--sb-score-width)',
                    left: 'calc(var(--sb-rank-width) + var(--sb-team-width))',
                  }}
                >
                  {t('game.scoreboardTable.headers.score')}
                </th>

                {/* 题目分类列 */}
                {categoryEntries.map(([category, categoryIchallenges]) => (
                  <th
                    key={category}
                    scope="colgroup"
                    className="text-center text-neutral-400 font-mono text-sm border-r border-neutral-300/30"
                    colSpan={categoryIchallenges.length}
                  >
                    <div className="flex items-center justify-center gap-2 p-3">
                      <span className="text-geek-400">#</span>
                      <span className="block max-w-full truncate">{category}</span>
                      <span className="text-xs text-neutral-500">({categoryIchallenges.length})</span>
                    </div>
                  </th>
                ))}
              </tr>

              {/* 题目名称行 */}
              <tr className="bg-neutral-800/40 border-t border-neutral-600/40">
                <th
                  className="sticky left-0 z-10 bg-neutral-800/60 p-2 border-r border-neutral-600/40"
                  style={{ width: 'var(--sb-rank-width)', minWidth: 'var(--sb-rank-width)' }}
                ></th>
                <th
                  className="sticky z-10 bg-neutral-800/60 p-2 border-r border-neutral-600/40"
                  style={{
                    width: 'var(--sb-team-width)',
                    minWidth: 'var(--sb-team-width)',
                    left: 'var(--sb-rank-width)',
                  }}
                ></th>
                <th
                  className="sticky z-10 bg-neutral-800/60 p-2 border-r border-neutral-600/40"
                  style={{
                    width: 'var(--sb-score-width)',
                    minWidth: 'var(--sb-score-width)',
                    left: 'calc(var(--sb-rank-width) + var(--sb-team-width))',
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
                        <span className="block max-w-full truncate">{challenge.name}</span>
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
                  className="border-t border-neutral-300/10 hover:bg-neutral-300/5"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: teamIndex * 0.05 }}
                >
                  {/* 排名 */}
                  <td
                    className="sticky left-0 z-10 bg-neutral-800/70 p-3 text-center text-neutral-300 font-mono tabular-nums border-r border-neutral-600/40"
                    style={{ width: 'var(--sb-rank-width)', minWidth: 'var(--sb-rank-width)' }}
                  >
                    {getTeamRank(teamIndex)}
                  </td>

                  {/* 队伍信息 */}
                  <td
                    className="sticky z-10 bg-neutral-800/70 p-3 border-r border-neutral-600/40"
                    style={{
                      width: 'var(--sb-team-width)',
                      minWidth: 'var(--sb-team-width)',
                      left: 'var(--sb-rank-width)',
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
                        <div className="text-xs text-neutral-500 tabular-nums">
                          {t('game.scoreboardTable.members', { count: team.users })}
                        </div>
                      </div>
                    </div>
                  </td>

                  {/* 分数 */}
                  <td
                    className="sticky z-10 bg-neutral-800/70 p-3 text-center text-geek-400 font-mono tabular-nums border-r border-neutral-600/40"
                    style={{
                      width: 'var(--sb-score-width)',
                      minWidth: 'var(--sb-score-width)',
                      left: 'calc(var(--sb-rank-width) + var(--sb-team-width))',
                    }}
                  >
                    {team.score}
                  </td>

                  {/* 题目状态 */}
                  {challengeColumns.map((challenge) => {
                    const solved = teamChallengeIndexes[teamIndex].get(challenge.id)?.solved > 0;

                    return (
                      <td
                        key={challenge.id}
                        className="group/cell p-3 text-center border-r border-neutral-300/30"
                        style={{ width: '80px' }}
                      >
                        <div
                          className="transition-transform duration-200 group-hover/cell:scale-110 flex items-center justify-center font-mono text-lg"
                          title={`${team.name} - ${challenge.name}: ${
                            solved ? t('game.scoreboardTable.status.solved') : t('game.scoreboardTable.status.unsolved')
                          }`}
                        >
                          {solved ? (
                            <span className="text-green-500">✓</span>
                          ) : (
                            <span className="text-neutral-600">•</span>
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

      {/* 分页 */}
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

export default ScoreboardTable;
