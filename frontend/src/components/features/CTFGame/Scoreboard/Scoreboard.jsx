/**
 * 积分榜主组件
 * @param {Object} props
 * @param {Object} props.stats - 总览统计数据
 * @param {number} props.stats.totalTeams - 总队伍数
 * @param {number} props.stats.totalSolves - 总解题数
 * @param {number} props.stats.highestScore - 最高分数
 * @param {string} props.stats.lastSubmitTime - 最后提交时间 (例: "2 mins ago")
 * @param {Array<Object>} props.teams - 队伍列表
 * @param {number} props.teams[].rank - 队伍排名
 * @param {string} props.teams[].name - 队伍名称
 * @param {string} props.teams[].picture - 队伍头像URL
 * @param {number} props.teams[].score - 队伍分数
 * @param {number} props.teams[].solved - 已解题目数量
 * @param {string} props.teams[].lastSubmit - 最后提交时间 (例: "2024-03-15 14:30:22")
 * @param {Object} props.teams[].solves - 各类型题目解题数
 * @param {number} props.teams[].solves.WEB - Web类题目解题数
 * @param {number} props.teams[].solves.CRYPTO - 密码学类题目解题数
 * @param {number} props.teams[].solves.PWN - PWN类题目解题数
 * @param {number} props.teams[].solves.REVERSE - 逆向类题目解题数
 * @param {number} props.teams[].solves.MISC - 杂项类题目解题数
 * @param {Object|null} props.userTeam - 当前用户所在队伍信息, 格式同teams元素, null表示未加入队伍
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.totalPages - 总页数
 * @param {Function} props.onPageChange - 页码变化处理函数
 * @param {string} props.viewMode - 视图模式 ('ranking' | 'table')
 * @param {Function} props.onViewModeChange - 视图切换回调
 * @param {Array} props.challenges - 题目列表（表格视图用）
 * @param {number} props.totalCount - 总队伍数（表格视图用）
 * @param {number} props.tableCurrentPage - 表格视图当前页码
 * @param {number} props.tableTotalCount - 表格视图总条目数
 * @param {Function} props.onTablePageChange - 表格视图页码变化处理函数
 * @param {boolean} props.tableLoading - 表格视图加载状态
 */

import ScoreboardRanking from './ScoreboardRanking.jsx';
import ScoreboardTable from './ScoreboardTable';
import ScoreboardTimeline from './ScoreboardTimeline.jsx';
import { Pagination } from '../../../../components/common';
import { useTranslation } from 'react-i18next';

function Scoreboard({
  teams,
  currentPage,
  totalPages,
  onPageChange,
  viewMode = 'ranking',
  challenges = [],
  totalCount = 0,
  tableCurrentPage,
  onTablePageChange,
}) {
  const { t, i18n } = useTranslation();

  return (
    <div className="contest-container mx-auto space-y-6">
      {/* 视图内容 */}
      {viewMode === 'ranking' ? (
        <>
          <ScoreboardRanking
            teams={teams}
            locale={i18n.language || 'en-US'}
            labels={{
              rank: t('game.scoreboard.headers.rank'),
              team: t('game.scoreboard.headers.team'),
              score: t('game.scoreboard.headers.score'),
              challenges: t('game.scoreboard.headers.challenges'),
              lastSubmit: t('game.scoreboard.headers.lastSubmit'),
              total: t('game.scoreboard.total'),
            }}
            emptyMessage={t('common.noData')}
            footer={
              totalPages > 1 ? (
                <div className="mt-6 flex justify-center">
                  <Pagination
                    current={currentPage}
                    total={totalPages}
                    onChange={onPageChange}
                    showTotal={true}
                    totalItems={totalCount}
                    animate
                  />
                </div>
              ) : null
            }
          />
        </>
      ) : viewMode === 'table' ? (
        /* 表格视图 */
        <ScoreboardTable
          teams={teams}
          challenges={challenges}
          totalCount={totalCount}
          currentPage={tableCurrentPage || 1}
          pageSize={10}
          onPageChange={onTablePageChange}
          isAdmin={false}
          PaginationComponent={Pagination}
        />
      ) : (
        <ScoreboardTimeline timelineData={challenges} />
      )}
    </div>
  );
}

export default Scoreboard;
