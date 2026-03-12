/**
 * 管理后台积分榜组件
 * @param {Object} props
 * @param {Array<Object>} props.teams - 队伍列表
 * @param {number} props.teams[].rank - 队伍排名
 * @param {string} props.teams[].name - 队伍名称
 * @param {string} props.teams[].picture - 队伍头像URL
 * @param {number} props.teams[].score - 队伍分数
 * @param {number} props.teams[].solved - 已解题目数量
 * @param {string} props.teams[].lastSubmit - 最后提交时间 (例: "2024-03-15 14:30:22")
 * @param {number} props.currentPage - 当前页码
 * @param {number} props.pageSize - 每页显示数量
 * @param {number} props.totalCount - 总队伍数
 * @param {Function} props.onPageChange - 页码变化处理函数
 * @param {Function} props.onRowClick - 行点击回调
 */

import { useTranslation } from 'react-i18next';
import { Pagination } from '../../../../components/common';
import ScoreboardRanking from '../../CTFGame/Scoreboard/ScoreboardRanking.jsx';

function AdminScoreboard({
  teams = [],
  currentPage = 1,
  pageSize = 6,
  totalCount = 0,
  onPageChange,
  onRowClick,
}) {
  const { t, i18n } = useTranslation();

  return (
    <ScoreboardRanking
      teams={teams}
      locale={i18n.language || 'en-US'}
      labels={{
        rank: t('admin.contests.scoreboard.headers.rank'),
        team: t('admin.contests.scoreboard.headers.team'),
        score: t('admin.contests.scoreboard.headers.score'),
        challenges: t('admin.contests.scoreboard.headers.challenges'),
        lastSubmit: t('admin.contests.scoreboard.headers.lastSubmit'),
        total: t('admin.contests.scoreboard.total'),
      }}
      emptyMessage={t('common.noData')}
      onRowClick={onRowClick}
      footer={
        totalCount > pageSize ? (
          <div className="mt-6">
            <Pagination
              total={Math.ceil(totalCount / pageSize)}
              current={currentPage}
              pageSize={pageSize}
              onChange={onPageChange}
              showTotal
              totalItems={totalCount}
            />
          </div>
        ) : null
      }
    />
  );
}

export default AdminScoreboard;
