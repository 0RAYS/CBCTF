import { useTranslation } from 'react-i18next';
import { Pagination } from '../../common';
import LeaderboardRanking from './LeaderboardRanking.jsx';

/** 管理后台积分榜组件 */
function AdminRanking({
  teams = [],
  currentPage = 1,
  pageSize = 6,
  totalCount = 0,
  onPageChange,
  onRowClick,
  loading,
}) {
  const { t, i18n } = useTranslation();

  return (
    <LeaderboardRanking
      teams={teams}
      locale={i18n.language || 'en-US'}
      labels={{
        rank: t('common.leaderboard.headers.rank'),
        team: t('common.leaderboard.headers.team'),
        score: t('common.leaderboard.headers.score'),
        challenges: t('common.leaderboard.headers.challenges'),
        lastSubmit: t('common.leaderboard.headers.lastSubmit'),
        total: t('common.leaderboard.total'),
      }}
      emptyMessage={t('common.noData')}
      loading={loading}
      onRowClick={onRowClick}
      footer={
        totalCount > pageSize ? (
          <div className="mt-6">
            <Pagination
              total={Math.ceil(totalCount / pageSize)}
              current={currentPage}
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

export default AdminRanking;
