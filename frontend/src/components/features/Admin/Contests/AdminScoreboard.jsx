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
import AdminTeamDetailDialog from './AdminTeamDetailDialog';

function AdminScoreboard({
  teams = [],
  currentPage = 1,
  pageSize = 6,
  totalCount = 0,
  onPageChange,
  onRowClick,
  showDetailDialog = false,
  detailTeam = null,
  detailTab = 'info',
  detailMembers = [],
  detailMembersLoading = false,
  detailSubmissions = [],
  detailSubmissionCount = 0,
  detailSubmissionPage = 1,
  detailWriteups = [],
  detailWriteupCount = 0,
  detailWriteupPage = 1,
  detailContainers = [],
  detailContainerCount = 0,
  detailContainerPage = 1,
  detailLoading = { submissions: false, writeups: false, traffic: false },
  onDetailClose,
  onDetailTabChange,
  onDetailPageChange,
  onDetailDownloadTraffic,
  onDetailDownloadWriteup,
  detailFlags = [],
  detailFlagsLoading = false,
}) {
  const { t, i18n } = useTranslation();

  return (
    <>
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

      <AdminTeamDetailDialog
        isOpen={showDetailDialog}
        onClose={onDetailClose}
        team={detailTeam}
        activeTab={detailTab}
        onTabChange={onDetailTabChange}
        members={detailMembers}
        membersLoading={detailMembersLoading}
        detailSubmissions={detailSubmissions}
        detailSubmissionCount={detailSubmissionCount}
        detailSubmissionPage={detailSubmissionPage}
        detailWriteups={detailWriteups}
        detailWriteupCount={detailWriteupCount}
        detailWriteupPage={detailWriteupPage}
        detailContainers={detailContainers}
        detailContainerCount={detailContainerCount}
        detailContainerPage={detailContainerPage}
        detailLoading={detailLoading}
        onDetailPageChange={onDetailPageChange}
        onDetailDownloadTraffic={onDetailDownloadTraffic}
        onDetailDownloadWriteup={onDetailDownloadWriteup}
        detailFlags={detailFlags}
        detailFlagsLoading={detailFlagsLoading}
      />
    </>
  );
}

export default AdminScoreboard;
